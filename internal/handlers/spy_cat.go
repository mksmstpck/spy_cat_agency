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

type spyCat struct {
	config   config.Config
	services *services.Services
}

func newSpyCat(
	config config.Config,
	services *services.Services,
) *spyCat {
	return &spyCat{
		config:   config,
		services: services,
	}
}

type spyCatInput struct {
	Name     string  `json:"name" binding:"required"`
	Breed    string  `json:"breed" binding:"required"`
	ExpYears int     `json:"years_experience"`
	Salary   float32 `json:"salary"`
}

func (input *spyCatInput) Validate() error {
	if strings.TrimSpace(input.Name) == "" {
		return &ValidationError{Field: "name", Message: "cat name cannot be empty"}
	}
	if strings.TrimSpace(input.Breed) == "" {
		return &ValidationError{Field: "breed", Message: "breed cannot be empty"}
	}
	if input.ExpYears < 0 {
		return &ValidationError{Field: "years_experience", Message: "experience cannot be negative"}
	}
	if input.Salary < 0 {
		return &ValidationError{Field: "salary", Message: "salary cannot be negative"}
	}
	return nil
}

func (h *spyCat) Create(c *gin.Context) {
	var catCreate spyCatInput

	if err := c.ShouldBindJSON(&catCreate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := catCreate.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	breed, err := h.services.Breed.GetByName(c.Request.Context(), strings.TrimSpace(catCreate.Breed))
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid breed",
			"details": err.Error(),
		})
		return
	}

	if breed == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Breed not found",
		})
		return
	}

	cat := &models.SpyCat{
		Name:     strings.TrimSpace(catCreate.Name),
		Breed:    *breed,
		ExpYears: catCreate.ExpYears,
		Salary:   catCreate.Salary,
	}

	if err := cat.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Model validation failed",
			"details": err.Error(),
		})
		return
	}

	cat, err = h.services.SpyCat.Create(c.Request.Context(), *cat)
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
			"error":   "Failed to create cat",
			"details": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func (h *spyCat) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cat ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cat ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	cat, err := h.services.SpyCat.GetByID(c.Request.Context(), newID)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve cat",
		})
		return
	}

	if cat == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Cat not found",
		})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (h *spyCat) GetAll(c *gin.Context) {
	cats, err := h.services.SpyCat.GetAll(c.Request.Context())
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve cats",
		})
		return
	}

	c.JSON(http.StatusOK, cats)
}

type spyCatSalaryUpdate struct {
	Salary *float32 `json:"salary" binding:"required"`
}

func (input *spyCatSalaryUpdate) Validate() error {
	if input.Salary == nil {
		return &ValidationError{Field: "salary", Message: "salary is required"}
	}
	if *input.Salary < 0 {
		return &ValidationError{Field: "salary", Message: "salary cannot be negative"}
	}
	return nil
}

func (h *spyCat) UpdateSalary(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cat ID is required",
		})
		return
	}

	catID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cat ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var catUpdate spyCatSalaryUpdate
	if err := c.ShouldBindJSON(&catUpdate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": "salary field is required and must be a number",
		})
		return
	}

	if err := catUpdate.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	err = h.services.SpyCat.UpdateSalary(c.Request.Context(), catID, *catUpdate.Salary)
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
			"error": "Failed to update cat salary",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type spyCatExpUpdate struct {
	ExpYears *int `json:"years_experience" binding:"required"`
}

func (input *spyCatExpUpdate) Validate() error {
	if input.ExpYears == nil {
		return &ValidationError{Field: "years_experience", Message: "years_experience is required"}
	}
	if *input.ExpYears < 0 {
		return &ValidationError{Field: "years_experience", Message: "experience cannot be negative"}
	}
	return nil
}

func (h *spyCat) UpdateExpYears(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cat ID is required",
		})
		return
	}

	catID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cat ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	var catUpdate spyCatExpUpdate
	if err := c.ShouldBindJSON(&catUpdate); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": "years_experience field is required and must be a number",
		})
		return
	}

	if err := catUpdate.Validate(); err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	err = h.services.SpyCat.UpdateExperience(c.Request.Context(), catID, *catUpdate.ExpYears)
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
			"error": "Failed to update cat experience",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *spyCat) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cat ID is required",
		})
		return
	}

	newID, err := uuid.Parse(id)
	if err != nil {
		logrus.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid cat ID format",
			"details": "ID must be a valid UUID",
		})
		return
	}

	err = h.services.SpyCat.Delete(c.Request.Context(), newID)
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
			"error": "Failed to delete cat",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
