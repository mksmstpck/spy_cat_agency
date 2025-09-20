package models

import (
	"time"

	"github.com/google/uuid"
)

type Target struct {
	ID        uuid.UUID
	MissionID uuid.UUID
	Name      string
	Country   string
	Notes     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
