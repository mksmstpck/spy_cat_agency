package services

import "github.com/mksmstpck/spy_cat_agency/internal/db"

type Services struct {
	Breed   breed
	SpyCat  spyCat
	Mission mission
	Target  target
}

func NewServices(db db.DB) *Services {
	return &Services{
		Breed:   *newBreed(db),
		SpyCat:  *newSpyCat(db),
		Mission: *newMission(db),
		Target:  *newTarget(db),
	}
}
