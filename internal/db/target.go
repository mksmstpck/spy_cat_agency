package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/sirupsen/logrus"
)

type target struct {
	conn *pgxpool.Pool
}

func newTarget(conn *pgxpool.Pool) *target {
	return &target{
		conn: conn,
	}
}

func (db *target) Create(ctx context.Context, target models.Target) (*models.Target, error) {
	err := db.conn.QueryRow(
		ctx,
		`INSERT INTO targets (mission_id, name, country, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, completed, created_at, updated_at`,
		target.MissionID,
		target.Name,
		target.Country,
		target.Notes,
	).Scan(&target.ID, &target.Completed, &target.CreatedAt, &target.UpdatedAt)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &target, nil
}

func (db *target) GetByID(ctx context.Context, id uuid.UUID) (*models.Target, error) {
	var target models.Target
	err := db.conn.QueryRow(
		ctx,
		`SELECT id, mission_id, name, country, notes, completed, created_at, updated_at
		FROM targets
		WHERE id = $1`,
		id,
	).Scan(
		&target.ID,
		&target.MissionID,
		&target.Name,
		&target.Country,
		&target.Notes,
		&target.Completed,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.Error(err)
		return nil, err
	}
	return &target, nil
}

func (db *target) UpdateCompleted(ctx context.Context, id uuid.UUID, completed bool) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE targets
		SET completed = $1
		WHERE id = $2`,
		completed,
		id,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (db *target) UpdateNotes(ctx context.Context, id uuid.UUID, notes string) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE targets
		SET notes = $1
		WHERE id = $2`,
		notes,
		id,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (db *target) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.conn.Exec(
		ctx,
		"DELETE FROM targets WHERE id = $1",
		id,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
