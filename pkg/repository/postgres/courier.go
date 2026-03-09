package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type PostgresCourierRepo struct {
	courierDB *sql.DB
}

func NewPostgresCourierRepo(courierDB *sql.DB) *PostgresCourierRepo {
	return &PostgresCourierRepo{courierDB: courierDB}
}

func EnsureCourierExists(ctx context.Context, db *sql.DB, courierID uuid.UUID) error {
	const query = "SELECT id FROM COURIERS WHERE id = $1"
	var id uuid.UUID
	if err := db.QueryRowContext(ctx, query, courierID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("courier not found")
		}
		return err
	}
	return nil
}
