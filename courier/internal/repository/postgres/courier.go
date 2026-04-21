package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kabanya/YAFDS/courier/internal/domain"
	"github.com/google/uuid"
)

type courierPostgresRepository struct {
	db *sql.DB
}

func NewCourierPostgresRepository(db *sql.DB) domain.CourierRepo {
	return &courierPostgresRepository{db: db}
}

func (r *courierPostgresRepository) ListCouriers(ctx context.Context) ([]domain.Courier, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, transport_type, is_active FROM COURIERS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []domain.Courier
	for rows.Next() {
		var c domain.Courier
		if err := rows.Scan(&c.ID, &c.Name, &c.TransportType, &c.IsActive); err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}
	return couriers, nil
}

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
