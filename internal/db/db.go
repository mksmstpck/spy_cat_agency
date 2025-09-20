package db

import "github.com/jackc/pgx/v5/pgxpool"

type DB struct {
	Breed   breed
	SpyCat  spyCat
	Mission mission
	Target  target
}

func NewDB(conn *pgxpool.Pool) *DB {
	return &DB{
		Breed:   *newBreed(conn),
		SpyCat:  *newSpyCat(conn),
		Mission: *newMission(conn),
		Target:  *newTarget(conn),
	}
}
