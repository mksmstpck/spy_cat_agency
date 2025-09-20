package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mksmstpck/spy_cat_agency/internal/config"
	"github.com/mksmstpck/spy_cat_agency/internal/services"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	spyCat *spyCat
	config config.Config
}

func NewHandlers(
	config config.Config,
	services *services.Services,
) *Handlers {
	return &Handlers{
		spyCat: newSpyCat(config, services),
		config: config,
	}
}

func (h *Handlers) HandleAll(ctx context.Context) {
	r := gin.Default()

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
		cat.PUT("/salary", h.spyCat.UpdateSalary)
		cat.PUT("/experience", h.spyCat.UpdateExpYears)
		cat.DELETE("/:id", h.spyCat.Delete)
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
