package services

import "github.com/mksmstpck/spy_cat_agency/internal/db"

type Services struct {
	Breed  breed
	SpyCat spyCat
}

func NewServices(db db.DB) *Services {
	return &Services{
		Breed:  *newBreed(db),
		SpyCat: *newSpyCat(db),
	}
}
