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

type target struct {
	config   config.Config
	services *services.Services
}

func newTarget(
	config config.Config,
	services *services.Services,
) *target {
	return &target{
		config:   config,
		services: services,
	}
}

type targetCreate struct {
	MissionID uuid.UUID `json:"mission_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Country   string    `json:"country" binding:"required"`
	Notes     string    `json:"notes"`
}

func (input *targetCreate) Validate() error {
	if strings.TrimSpace(input.Name) == "" {
		return &ValidationError{Field: "name", Message: "target name cannot be empty"}
	}
	if strings.TrimSpace(input.Country) == "" {
		return &ValidationError{Field: "country", Message: "target country cannot be empty"}
	}
	return nil
}

func (h *target) Create(c *gin.Context) {
	var targetCreate targetCreate

	if err := c.ShouldBindJSON(&targetCreate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := targetCreate.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	target := models.Target{
		MissionID: targetCreate.MissionID,
		Name:      strings.TrimSpace(targetCreate.Name),
		Country:   strings.TrimSpace(targetCreate.Country),
		Notes:     targetCreate.Notes,
	}

	createdTarget, err := h.services.Target.Create(c.Request.Context(), target)
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
			"error": "Failed to create target",
		})
		return
	}

	c.JSON(http.StatusCreated, createdTarget)
}

func (h *target) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Target ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid target ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	target, err := h.services.Target.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve target",
		})
		return
	}

	if target == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Target not found",
		})
		return
	}

	c.JSON(http.StatusOK, target)
}

type targetUpdateCompleted struct {
	Completed *bool `json:"completed" binding:"required"`
}

func (h *target) UpdateCompleted(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Target ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid target ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var targetUpdate targetUpdateCompleted
	if err := c.ShouldBindJSON(&targetUpdate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": "completed field is required and must be a boolean",
		})
		return
	}

	err = h.services.Target.UpdateCompleted(c.Request.Context(), newID, *targetUpdate.Completed)
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
			"error": "Failed to update target",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type targetUpdateNotes struct {
	Notes *string `json:"notes" binding:"required"`
}

func (h *target) UpdateNotes(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Target ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid target ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var targetUpdate targetUpdateNotes
	if err := c.ShouldBindJSON(&targetUpdate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": "notes field is required",
		})
		return
	}

	err = h.services.Target.UpdateNotes(c.Request.Context(), newID, *targetUpdate.Notes)
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
			"error": "Failed to update target notes",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *target) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Target ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid target ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	err = h.services.Target.Delete(c.Request.Context(), newID)
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
			"error": "Failed to delete target",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
