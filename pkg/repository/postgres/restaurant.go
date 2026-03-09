package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type restaurantPostgresRepo struct {
	db *sql.DB
}

func NewRestaurantPostgresRepo(db *sql.DB) repositoryModels.RestaurantRepo {
	return &restaurantPostgresRepo{db: db}
}

func EnsureRestaurantExists(ctx context.Context, db *sql.DB, restaurantID string) error {
	const query = "SELECT id FROM RESTAURANTS WHERE id = $1"
	var id string
	if err := db.QueryRowContext(ctx, query, restaurantID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("restaurant not found")
		}
		return err
	}
	return nil
}

func (r *restaurantPostgresRepo) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, address, status FROM RESTAURANTS")
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
