package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	spyCat *spyCat
}

func NewHandlers(spyCat *spyCat) *Handlers {
	return &Handlers{
		spyCat: spyCat,
	}
}

func (h *Handlers) HandleAll() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "https://accounts.google.com"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	cat := r.Group("cat")
	{
		cat.GET("/", h.spyCat.GetAll)
		cat.POST("/", h.spyCat.Create)
		cat.PUT("/salary", h.spyCat.UpdateSalary)
		cat.PUT("/experience", h.spyCat.UpdateExpYears)
		cat.DELETE("/:id", h.spyCat.Delete)
	}
}
