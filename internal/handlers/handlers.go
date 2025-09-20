package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mksmstpck/spy_cat_agency/internal/config"
	"github.com/mksmstpck/spy_cat_agency/internal/services"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	spyCat  *spyCat
	mission *mission
	target  *target
	config  config.Config
}

func NewHandlers(
	config config.Config,
	services *services.Services,
) *Handlers {
	return &Handlers{
		spyCat:  newSpyCat(config, services),
		mission: newMission(config, services),
		target:  newTarget(config, services),
		config:  config,
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Index   *int   `json:"index,omitempty"`
}

func (e *ValidationError) Error() string {
	if e.Index != nil {
		return e.Field + "[" + string(rune(*e.Index)) + "]: " + e.Message
	}
	return e.Field + ": " + e.Message
}

func isBusinessLogicError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	businessLogicKeywords := []string{
		"cannot delete",
		"cannot add target",
		"cannot update notes",
		"must have at least",
		"cannot have more than",
		"mission completed",
		"targets still incomplete",
		"assigned to cat",
		"salary must be",
		"experience cannot be",
	}

	for _, keyword := range businessLogicKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

func (h *Handlers) HandleAll(ctx context.Context) {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "https://accounts.google.com"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(RequestLogger())

	cat := r.Group("cat")
	{
		cat.GET("/", h.spyCat.GetAll)
		cat.GET("/:id", h.spyCat.GetByID)
		cat.POST("/", h.spyCat.Create)
		cat.PUT("/salary/:id", h.spyCat.UpdateSalary)
		cat.PUT("/experience/:id", h.spyCat.UpdateExpYears)
		cat.DELETE("/:id", h.spyCat.Delete)
	}

	mission := r.Group("mission")
	{
		mission.GET("/", h.mission.GetAll)
		mission.GET("/:id", h.mission.GetByID)
		mission.POST("/", h.mission.Create)
		mission.PUT("/:id/completed", h.mission.UpdateCompleted)
		mission.PUT("/:id/assign", h.mission.AssignCat)
		mission.DELETE("/:id", h.mission.Delete)
	}

	target := r.Group("target")
	{
		target.GET("/:id", h.target.GetByID)
		target.POST("/", h.target.Create)
		target.PUT("/:id/completed", h.target.UpdateCompleted)
		target.PUT("/:id/notes", h.target.UpdateNotes)
		target.DELETE("/:id", h.target.Delete)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", h.config.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Gin server error: %s", err)
		}
	}()

	<-ctx.Done()

	logrus.Info("Shutting down Gin server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("Gin server forced to shut down: %s", err)
	}
}
