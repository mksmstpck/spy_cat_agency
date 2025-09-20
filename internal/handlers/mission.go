package handlers

import (
	"net/http"

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
	Title         string        `json:"title"`
	Description   *string       `json:"description,omitempty"`
	AssignedCatID *uuid.UUID    `json:"assigned_cat_id,omitempty"`
	Targets       []targetInput `json:"targets"`
}

type targetInput struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Notes   string `json:"notes"`
}

func (h *mission) Create(c *gin.Context) {
	var missionCreate missionInput

	err := c.Bind(&missionCreate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mission := models.Mission{
		Title:         missionCreate.Title,
		Description:   missionCreate.Description,
		AssignedCatID: missionCreate.AssignedCatID,
	}

	targets := make([]models.Target, len(missionCreate.Targets))
	for i, targetInput := range missionCreate.Targets {
		targets[i] = models.Target{
			Name:    targetInput.Name,
			Country: targetInput.Country,
			Notes:   targetInput.Notes,
		}
	}

	createdMission, err := h.services.Mission.Create(c.Request.Context(), mission, targets)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMission)
}

func (h *mission) GetByID(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mission, err := h.services.Mission.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mission)
}

func (h *mission) GetAll(c *gin.Context) {
	missions, err := h.services.Mission.GetAll(c.Request.Context())
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, missions)
}

type missionUpdate struct {
	Completed *bool `json:"completed,omitempty"`
}

func (h *mission) UpdateCompleted(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var missionUpdate missionUpdate
	err = c.Bind(&missionUpdate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if missionUpdate.Completed == nil {
		logrus.Error("completed field is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "completed field is required"})
		return
	}

	err = h.services.Mission.UpdateCompleted(c.Request.Context(), newID, *missionUpdate.Completed)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type assignCatInput struct {
	CatID *uuid.UUID `json:"cat_id"`
}

func (h *mission) AssignCat(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var assignInput assignCatInput
	err = c.Bind(&assignInput)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Mission.UpdateAssignedCat(c.Request.Context(), newID, assignInput.CatID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *mission) Delete(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Mission.Delete(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
