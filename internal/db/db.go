package db

import "github.com/jackc/pgx/v5/pgxpool"

type DB struct {
	Breed  breed
	SpyCat spyCat
}

func NewDB(conn *pgxpool.Pool) *DB {
	return &DB{
		Breed:  *newBreed(conn),
		SpyCat: *newSpyCat(conn),
	}
}
