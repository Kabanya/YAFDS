package repository

import (
	"context"
	"database/sql"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type restaurantPostgresRepository struct {
	db *sql.DB
}

func NewRestaurantPostgresRepository(db *sql.DB) repositoryModels.RestaurantRepo {
	return &restaurantPostgresRepository{db: db}
}

func (r *restaurantPostgresRepository) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT emp_id, name, address, status FROM RESTAURANTS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var res models.Restaurant
		if err := rows.Scan(&res.ID, &res.Name, &res.Address, &res.Status); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, res)
	}
	return restaurants, nil
}
