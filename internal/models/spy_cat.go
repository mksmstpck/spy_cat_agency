package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type SpyCat struct {
	ID        uuid.UUID
	Name      string
	ExpYears  int
	Breed     Breed
	Salary    float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s SpyCat) Validate() error {
	if s.Name == "" {
		return errors.New("cat name cannot be empty")
	}
	if s.Salary < 0 {
		return errors.New("salary cannot be negative")
	}

	return nil
}
