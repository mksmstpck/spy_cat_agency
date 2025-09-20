package models

import (
	"time"

	"github.com/google/uuid"
)

type Mission struct {
	ID            uuid.UUID
	Title         string
	Description   *string
	AssignedCatID *uuid.UUID
	Targets       []Target
	Completed     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
