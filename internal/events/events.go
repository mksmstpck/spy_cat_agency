package events

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mksmstpck/spy_cat_agency/internal/config"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/mksmstpck/spy_cat_agency/internal/services"
	"github.com/sirupsen/logrus"
)

type Events struct {
	services services.Services
	config   config.Config
}

func NewEvents(services services.Services, config config.Config) *Events {
	return &Events{
		services: services,
		config:   config,
	}
}

type breed struct {
	Name  string `json:"name"`
	ApiID string `json:"id"`
}

func (e *Events) LoadBreeds(ctx context.Context) error {
	resp, err := http.Get(e.config.TheCatApiUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return err
	}

	var breeds []breed

	if err := json.Unmarshal(body, &breeds); err != nil {
		logrus.Error(err)
		return err
	}

	mBreed := models.Breed{}
	for _, breed := range breeds {
		mBreed.ApiID = breed.ApiID
		mBreed.Name = breed.Name

		_, err := e.services.Breed.Create(ctx, mBreed)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}

	return nil
}
