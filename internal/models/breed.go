package models

import (
	"time"

	"github.com/google/uuid"
)

type Breed struct {
	ID        uuid.UUID
	Name      string
	ApiID     string
	CreatedAt time.Time
}
