package services

import (
	"context"

	"github.com/mksmstpck/spy_cat_agency/internal/db"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
)

type breed struct {
	db db.DB
}

func newBreed(db db.DB) *breed {
	return &breed{
		db: db,
	}
}

func (s *breed) Create(ctx context.Context, breed models.Breed) (*models.Breed, error) {
	return s.db.Breed.Create(ctx, breed)
}

func (s *breed) GetAll(ctx context.Context) ([]models.Breed, error) {
	return s.db.Breed.GetAll(ctx)
}

func (s *breed) GetByName(ctx context.Context, name string) (*models.Breed, error) {
	return s.db.Breed.GetByName(ctx, name)
}
