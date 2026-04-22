package repository

import (
	"context"
	"database/sql"

	"github.com/Kabanya/YAFDS/courier/internal/domain"
)

type PostgresCourierRepo struct {
	courierDB *sql.DB
}
type courierPostgresRepository struct {
	db *sql.DB
}

func NewCourierPostgresRepository(db *sql.DB) domain.CourierRepo {
	return &courierPostgresRepository{db: db}
}

func NewPostgresCourierRepo(courierDB *sql.DB) *PostgresCourierRepo {
	return &PostgresCourierRepo{courierDB: courierDB}
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
