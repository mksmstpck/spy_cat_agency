package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mksmstpck/spy_cat_agency/internal/db"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
)

type target struct {
	db db.DB
}

func newTarget(db db.DB) *target {
	return &target{
		db: db,
	}
}

func (s *target) Create(ctx context.Context, target models.Target) (*models.Target, error) {
	if target.Name == "" {
		return nil, errors.New("target name cannot be empty")
	}
	if target.Country == "" {
		return nil, errors.New("target country cannot be empty")
	}

	return s.db.Target.Create(ctx, target)
}

func (s *target) GetByID(ctx context.Context, id uuid.UUID) (*models.Target, error) {
	return s.db.Target.GetByID(ctx, id)
}

func (s *target) UpdateCompleted(ctx context.Context, id uuid.UUID, completed bool) error {
	return s.db.Target.UpdateCompleted(ctx, id, completed)
}

func (s *target) UpdateNotes(ctx context.Context, id uuid.UUID, notes string) error {
	return s.db.Target.UpdateNotes(ctx, id, notes)
}

func (s *target) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.Target.Delete(ctx, id)
}
