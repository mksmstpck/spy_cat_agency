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

type spyCat struct {
	config   *config.Config
	services *services.Services
}

func newSpyCat(
	config *config.Config,
	services *services.Services,
) *spyCat {
	return &spyCat{
		config:   config,
		services: services,
	}
}

type spyCatInput struct {
	Name     string  `json:"name"`
	Breed    string  `json:"breed"`
	ExpYears int     `json:"years_experience"`
	Salary   float32 `json:"salary"`
}

func (h *spyCat) Create(c *gin.Context) {
	var catCreate spyCatInput

	err := c.Bind(&catCreate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	breed, err := h.services.Breed.GetByName(c.Request.Context(), catCreate.Breed)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	cat := &models.SpyCat{
		Name:   catCreate.Name,
		Breed:  *breed,
		Salary: catCreate.Salary,
	}

	cat, err = h.services.SpyCat.Create(c.Request.Context(), *cat)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *spyCat) GetByID(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	cat, err := h.services.SpyCat.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, cat)
}

func (h *spyCat) GetAll(c *gin.Context) {
	cats, err := h.services.SpyCat.GetAll(c.Request.Context())
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, cats)
}

type spyCatUpdate struct {
	ID       uuid.UUID `json:"id"`
	ExpYears int       `json:"exp_years"`
	Salary   float32   `json:"salary"`
}

func (h *spyCat) UpdateSalary(c *gin.Context) {
	var catUpdate spyCatUpdate

	err := c.Bind(&catUpdate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	cat := models.SpyCat{
		ID:       catUpdate.ID,
		ExpYears: catUpdate.ExpYears,
		Salary:   catUpdate.Salary,
	}

	err = h.services.SpyCat.UpdateSalary(c.Request.Context(), cat.ID, cat.Salary)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *spyCat) UpdateExpYears(c *gin.Context) {
	var catUpdate spyCatUpdate

	err := c.Bind(&catUpdate)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	cat := models.SpyCat{
		ID:       catUpdate.ID,
		ExpYears: catUpdate.ExpYears,
		Salary:   catUpdate.Salary,
	}

	err = h.services.SpyCat.UpdateExperience(c.Request.Context(), cat.ID, cat.ExpYears)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *spyCat) Delete(c *gin.Context) {
	id := c.Param("id")
	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	err = h.services.SpyCat.Delete(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusNoContent, nil)
}
