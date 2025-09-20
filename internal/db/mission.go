package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mksmstpck/spy_cat_agency/internal/models"
	"github.com/sirupsen/logrus"
)

type mission struct {
	conn *pgxpool.Pool
}

func newMission(conn *pgxpool.Pool) *mission {
	return &mission{
		conn: conn,
	}
}

func (db *mission) Create(ctx context.Context, mission models.Mission, targets []models.Target) (*models.Mission, error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		`INSERT INTO missions (title, description, assigned_cat_id)
		VALUES ($1, $2, $3)
		RETURNING id, completed, created_at, updated_at`,
		mission.Title,
		mission.Description,
		mission.AssignedCatID,
	).Scan(&mission.ID, &mission.Completed, &mission.CreatedAt, &mission.UpdatedAt)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	
	for i := range targets {
		targets[i].MissionID = mission.ID
		err = tx.QueryRow(
			ctx,
			`INSERT INTO targets (mission_id, name, country, notes)
			VALUES ($1, $2, $3, $4)
			RETURNING id, completed, created_at, updated_at`,
			targets[i].MissionID,
			targets[i].Name,
			targets[i].Country,
			targets[i].Notes,
		).Scan(&targets[i].ID, &targets[i].Completed, &targets[i].CreatedAt, &targets[i].UpdatedAt)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		logrus.Error(err)
		return nil, err
	}

	mission.Targets = targets
	return &mission, nil
}

func (db *mission) GetAll(ctx context.Context) ([]models.Mission, error) {
	rows, err := db.conn.Query(
		ctx,
		`SELECT id, title, description, assigned_cat_id, completed, created_at, updated_at
		FROM missions
		ORDER BY created_at DESC`,
	)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		err := rows.Scan(
			&mission.ID,
			&mission.Title,
			&mission.Description,
			&mission.AssignedCatID,
			&mission.Completed,
			&mission.CreatedAt,
			&mission.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		targets, err := db.getTargetsByMissionID(ctx, mission.ID)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		mission.Targets = targets

		missions = append(missions, mission)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return missions, nil
}

func (db *mission) GetByID(ctx context.Context, id uuid.UUID) (*models.Mission, error) {
	var mission models.Mission
	err := db.conn.QueryRow(
		ctx,
		`SELECT id, title, description, assigned_cat_id, completed, created_at, updated_at
		FROM missions
		WHERE id = $1`,
		id,
	).Scan(
		&mission.ID,
		&mission.Title,
		&mission.Description,
		&mission.AssignedCatID,
		&mission.Completed,
		&mission.CreatedAt,
		&mission.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logrus.Error(err)
		return nil, err
	}

	targets, err := db.getTargetsByMissionID(ctx, mission.ID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	mission.Targets = targets

	return &mission, nil
}

func (db *mission) UpdateCompleted(ctx context.Context, id uuid.UUID, completed bool) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE missions
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

func (db *mission) UpdateAssignedCat(ctx context.Context, id uuid.UUID, catID *uuid.UUID) error {
	_, err := db.conn.Exec(
		ctx,
		`UPDATE missions
		SET assigned_cat_id = $1
		WHERE id = $2`,
		catID,
		id,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (db *mission) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.conn.Exec(
		ctx,
		"DELETE FROM missions WHERE id = $1",
		id,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (db *mission) getTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]models.Target, error) {
	rows, err := db.conn.Query(
		ctx,
		`SELECT id, mission_id, name, country, notes, completed, created_at, updated_at
		FROM targets
		WHERE mission_id = $1
		ORDER BY created_at ASC`,
		missionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var target models.Target
		err := rows.Scan(
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
			return nil, err
		}
		targets = append(targets, target)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return targets, nil
}
