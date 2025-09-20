package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mksmstpck/spy_cat_agency/internal/db"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
)

type mission struct {
	db db.DB
}

func newMission(db db.DB) *mission {
	return &mission{
		db: db,
	}
}

func (s *mission) Create(ctx context.Context, mission models.Mission, targets []models.Target) (*models.Mission, error) {
	if len(targets) < 1 || len(targets) > 3 {
		return nil, errors.New("mission must have between 1 and 3 targets")
	}

	for _, target := range targets {
		if target.Name == "" {
			return nil, errors.New("target name cannot be empty")
		}
		if target.Country == "" {
			return nil, errors.New("target country cannot be empty")
		}
	}

	return s.db.Mission.Create(ctx, mission, targets)
}

func (s *mission) GetAll(ctx context.Context) ([]models.Mission, error) {
	return s.db.Mission.GetAll(ctx)
}

func (s *mission) GetByID(ctx context.Context, id uuid.UUID) (*models.Mission, error) {
	return s.db.Mission.GetByID(ctx, id)
}

func (s *mission) UpdateCompleted(ctx context.Context, id uuid.UUID, completed bool) error {
	return s.db.Mission.UpdateCompleted(ctx, id, completed)
}

func (s *mission) UpdateAssignedCat(ctx context.Context, id uuid.UUID, catID *uuid.UUID) error {
	return s.db.Mission.UpdateAssignedCat(ctx, id, catID)
}

func (s *mission) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.Mission.Delete(ctx, id)
}
