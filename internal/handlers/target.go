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
	MissionID uuid.UUID `json:"mission_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Notes     string    `json:"notes"`
}

func (h *target) Create(c *gin.Context) {
	var targetCreate targetCreate

	err := c.Bind(&targetCreate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target := models.Target{
		MissionID: targetCreate.MissionID,
		Name:      targetCreate.Name,
		Country:   targetCreate.Country,
		Notes:     targetCreate.Notes,
	}

	createdTarget, err := h.services.Target.Create(c.Request.Context(), target)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTarget)
}

func (h *target) GetByID(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := h.services.Target.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}

type targetUpdate struct {
	Completed *bool   `json:"completed,omitempty"`
	Notes     *string `json:"notes,omitempty"`
}

func (h *target) UpdateCompleted(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var targetUpdate targetUpdate
	err = c.Bind(&targetUpdate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if targetUpdate.Completed == nil {
		logrus.Error("completed field is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "completed field is required"})
		return
	}

	err = h.services.Target.UpdateCompleted(c.Request.Context(), newID, *targetUpdate.Completed)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *target) UpdateNotes(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var targetUpdate targetUpdate
	err = c.Bind(&targetUpdate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if targetUpdate.Notes == nil {
		logrus.Error("notes field is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "notes field is required"})
		return
	}

	err = h.services.Target.UpdateNotes(c.Request.Context(), newID, *targetUpdate.Notes)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *target) Delete(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.services.Target.Delete(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
