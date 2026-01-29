package repository

import (
	"context"
	"database/sql"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type courierPostgresRepository struct {
	db *sql.DB
}

func NewCourierPostgresRepository(db *sql.DB) repositoryModels.CourierRepo {
	return &courierPostgresRepository{db: db}
}

func (r *courierPostgresRepository) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT emp_id, name, transport_type, is_active FROM COURIERS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var couriers []models.Courier
	for rows.Next() {
		var c models.Courier
		if err := rows.Scan(&c.ID, &c.Name, &c.TransportType, &c.IsActive); err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}
	return couriers, nil
}
