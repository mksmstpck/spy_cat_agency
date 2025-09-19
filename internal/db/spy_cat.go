package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/sirupsen/logrus"
)

type spyCat struct {
	conn *pgxpool.Pool
}

func newSpyCat(conn *pgxpool.Pool) *spyCat {
	return &spyCat{
		conn: conn,
	}
}

func (db *spyCat) Create(ctx context.Context, cat models.SpyCat) (*models.SpyCat, error) {
	err := db.conn.QueryRow(
		ctx,
		`INSERT INTO cats (name, breed_id, years_experience, salary)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		cat.Name,
		cat.Breed.ID,
		cat.ExpYears,
		cat.Salary,
	).Scan(&cat.ID, &cat.CreatedAt, &cat.UpdatedAt)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &cat, nil
}

func (db *spyCat) GetAll(ctx context.Context) ([]models.SpyCat, error) {
	rows, err := db.conn.Query(
		ctx,
		`SELECT
			c.id,
			c.name,
			c.years_experience,
			c.salary,
			c.created_at,
			c.updated_at,
			b.id,
			b.api_id,
			b.name,
			b.created_at
		FROM cats c
		LEFT JOIN breeds b ON c.breed_id = b.id`,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	var cats []models.SpyCat

	for rows.Next() {
		var cat models.SpyCat
		var breed models.Breed

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.ExpYears,
			&cat.Salary,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&breed.ID,
			&breed.ApiID,
			&breed.Name,
			&breed.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		cat.Breed = breed
		cats = append(cats, cat)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return cats, nil
}

func (db *spyCat) GetByID(ctx context.Context, id uuid.UUID) (models.SpyCat, error) {
	var cat models.SpyCat

	err := db.conn.QueryRow(
		ctx,
		`SELECT c.id,
		c.name,
		c.years_experience,
		c.salary,
		c.created_at,
		c.updated_at,
		b.id,
		b.api_id,
		b.name,
		b.created_at
		 FROM cats c
		 LEFT JOIN breeds b ON sc.breed_id = b.id
		 WHERE sc.id = $1;`,
		id,
	).Scan(
		&cat.ID,
		&cat.Name,
		&cat.ExpYears,
		&cat.Salary,
		&cat.CreatedAt,
		&cat.UpdatedAt,
		&cat.Breed.ID,
		&cat.Breed.ApiID,
		&cat.Breed.Name,
		&cat.Breed.CreatedAt,
	)

	if err != nil {
		return models.SpyCat{}, err
	}

	return cat, nil
}

func (db *spyCat) UpdateSalary(ctx context.Context, id uuid.UUID, salary float32) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE cats
		SET salary = $1
		WHERE id = $2`,
		salary,
		id,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (db *spyCat) UpdateExperience(ctx context.Context, id uuid.UUID, exp int) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE cats
		SET years_experience = $1
		WHERE id = $2`,
		exp,
		id,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (db *spyCat) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.conn.Exec(
		ctx,
		"DELETE FROM cats WHERE id = $1",
		id.ID,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
