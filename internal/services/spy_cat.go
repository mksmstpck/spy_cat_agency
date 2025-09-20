package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mksmstpck/spy_cat_agency/internal/db"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
)

type spyCat struct {
	db db.DB
}

func newSpyCat(db db.DB) *spyCat {
	return &spyCat{
		db: db,
	}
}

func (s *spyCat) Create(ctx context.Context, cat models.SpyCat) (*models.SpyCat, error) {
	return s.db.SpyCat.Create(ctx, cat)
}

func (s *spyCat) GetAll(ctx context.Context) ([]models.SpyCat, error) {
	return s.db.SpyCat.GetAll(ctx)
}

func (s *spyCat) GetByID(ctx context.Context, id uuid.UUID) (*models.SpyCat, error) {
	return s.db.SpyCat.GetByID(ctx, id)
}

func (s *spyCat) UpdateSalary(ctx context.Context, id uuid.UUID, salary float32) error {
	if salary < 0 {
		return errors.New("salary must be >= 0")
	}
	return s.db.SpyCat.UpdateSalary(ctx, id, salary)
}

func (s *spyCat) UpdateExperience(ctx context.Context, id uuid.UUID, exp int) error {
	if exp < 0 {
		return errors.New("experience cannot be negative")
	}
	return s.db.SpyCat.UpdateExperience(ctx, id, exp)
}

func (s *spyCat) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.SpyCat.Delete(ctx, id)
}
