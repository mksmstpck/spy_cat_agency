package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mksmstpck/spy_cat_agency/internal/config"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/mksmstpck/spy_cat_agency/internal/services"
	"github.com/sirupsen/logrus"
)

type mission struct {
	config   config.Config
	services *services.Services
}

func newMission(
	config config.Config,
	services *services.Services,
) *mission {
	return &mission{
		config:   config,
		services: services,
	}
}

type missionInput struct {
	Title         string        `json:"title" binding:"required"`
	Description   *string       `json:"description,omitempty"`
	AssignedCatID *uuid.UUID    `json:"assigned_cat_id,omitempty"`
	Targets       []targetInput `json:"targets" binding:"required,min=1,max=3,dive"`
}

type targetInput struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
	Notes   string `json:"notes"`
}

func (input *missionInput) Validate() error {
	if strings.TrimSpace(input.Title) == "" {
		return &ValidationError{Field: "title", Message: "title cannot be empty"}
	}

	if len(input.Targets) < 1 || len(input.Targets) > 3 {
		return &ValidationError{Field: "targets", Message: "mission must have between 1 and 3 targets"}
	}

	for i, target := range input.Targets {
		if strings.TrimSpace(target.Name) == "" {
			return &ValidationError{Field: "targets", Message: "target name cannot be empty", Index: &i}
		}
		if strings.TrimSpace(target.Country) == "" {
			return &ValidationError{Field: "targets", Message: "target country cannot be empty", Index: &i}
		}
	}

	return nil
}

func (h *mission) Create(c *gin.Context) {
	var missionCreate missionInput

	if err := c.ShouldBindJSON(&missionCreate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := missionCreate.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	mission := models.Mission{
		Title:         strings.TrimSpace(missionCreate.Title),
		Description:   missionCreate.Description,
		AssignedCatID: missionCreate.AssignedCatID,
	}

	targets := make([]models.Target, len(missionCreate.Targets))
	for i, targetInput := range missionCreate.Targets {
		targets[i] = models.Target{
			Name:    strings.TrimSpace(targetInput.Name),
			Country: strings.TrimSpace(targetInput.Country),
			Notes:   targetInput.Notes,
		}
	}

	createdMission, err := h.services.Mission.Create(c.Request.Context(), mission, targets)
	if err != nil {
		logrus.Error(err)
		if isBusinessLogicError(err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error":   "Business rule violation",
				"details": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create mission",
			"details": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, createdMission)
}

func (h *mission) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Mission ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mission ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	mission, err := h.services.Mission.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve mission",
		})
		return
	}

	if mission == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Mission not found",
		})
		return
	}

	c.JSON(http.StatusOK, mission)
}

func (h *mission) GetAll(c *gin.Context) {
	missions, err := h.services.Mission.GetAll(c.Request.Context())
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve missions",
		})
		return
	}

	c.JSON(http.StatusOK, missions)
}

type missionUpdate struct {
	Completed *bool `json:"completed" binding:"required"`
}

func (h *mission) UpdateCompleted(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Mission ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mission ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var missionUpdate missionUpdate
	if err := c.ShouldBindJSON(&missionUpdate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": "completed field is required and must be a boolean",
		})
		return
	}

	err = h.services.Mission.UpdateCompleted(c.Request.Context(), newID, *missionUpdate.Completed)
	if err != nil {
		logrus.Error(err)
		if isBusinessLogicError(err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error":   "Business rule violation",
				"details": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update mission",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type assignCatInput struct {
	CatID *uuid.UUID `json:"cat_id"`
}

func (h *mission) AssignCat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Mission ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mission ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var assignInput assignCatInput
	if err := c.ShouldBindJSON(&assignInput); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.services.Mission.UpdateAssignedCat(c.Request.Context(), newID, assignInput.CatID)
	if err != nil {
		logrus.Error(err)
		if isBusinessLogicError(err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error":   "Business rule violation",
				"details": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to assign cat to mission",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *mission) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Mission ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid mission ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	err = h.services.Mission.Delete(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		if isBusinessLogicError(err) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error":   "Business rule violation",
				"details": err.Error(),
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete mission",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
