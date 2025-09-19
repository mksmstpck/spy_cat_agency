package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/sirupsen/logrus"
)

type breed struct {
	conn *pgxpool.Pool
}

func newBreed(conn *pgxpool.Pool) *breed {
	return &breed{
		conn: conn,
	}
}

func (db *breed) Create(ctx context.Context, breed models.Breed) (*models.Breed, error) {
	row := db.conn.QueryRow(
		ctx,
		`INSERT INTO breeds (api_id, name)
         VALUES ($1, $2)
         ON CONFLICT (api_id) DO NOTHING
         RETURNING id, created_at;`,
		breed.ApiID,
		breed.Name,
	)

	err := row.Scan(&breed.ID, &breed.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = db.conn.QueryRow(
				ctx,
				`SELECT id, created_at FROM breeds WHERE api_id = $1;`,
				breed.ApiID,
			).Scan(&breed.ID, &breed.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("breed exists but failed to fetch: %w", err)
			}
			return &breed, nil
		}
		return nil, fmt.Errorf("failed to insert breed: %w", err)
	}

	return &breed, nil
}

func (db *breed) GetAll(ctx context.Context) ([]models.Breed, error) {
	var breeds []models.Breed
	rows, err := db.conn.Query(
		ctx,
		`SELECT id, name, api_id, created_at FROM breeds`,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var breed models.Breed

		err := rows.Scan(
			&breed.ID,
			&breed.Name,
			&breed.ApiID,
			&breed.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		breeds = append(breeds, breed)
	}

	return breeds, nil
}

func (db *breed) GetByName(ctx context.Context, name string) (*models.Breed, error) {
	row := db.conn.QueryRow(
		ctx,
		`SELECT id, name, api_id, created_at FROM breeds WHERE name = $1;`,
		name,
	)

	var breed models.Breed

	err := row.Scan(
		&breed.ID,
		&breed.Name,
		&breed.ApiID,
		&breed.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &breed, nil
}
