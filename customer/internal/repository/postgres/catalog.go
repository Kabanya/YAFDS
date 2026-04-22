package repository

import (
	"context"
	"database/sql"

	customerrepo "github.com/Kabanya/YAFDS/customer/internal/repository"
	pkgErrors "github.com/Kabanya/YAFDS/pkg/errors"
	"github.com/google/uuid"
)

type catalogRepo struct {
	courierDB    *sql.DB
	restaurantDB *sql.DB
}

func NewCatalogRepo(courierDB, restaurantDB *sql.DB) customerrepo.CatalogRepository {
	return &catalogRepo{
		courierDB:    courierDB,
		restaurantDB: restaurantDB,
	}
}

func (r *catalogRepo) ListCouriers(ctx context.Context) ([]customerrepo.CourierCatalogDTO, error) {
	if r.courierDB == nil {
		return nil, pkgErrors.ErrRepositoryNotInitialized
	}

	rows, err := r.courierDB.QueryContext(ctx, "SELECT id, name, transport_type, is_active FROM COURIERS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	couriers := []customerrepo.CourierCatalogDTO{}
	for rows.Next() {
		var courier customerrepo.CourierCatalogDTO
		if err := rows.Scan(&courier.ID, &courier.Name, &courier.TransportType, &courier.IsActive); err != nil {
			return nil, err
		}
		couriers = append(couriers, courier)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return couriers, nil
}

func (r *catalogRepo) ListRestaurants(ctx context.Context) ([]customerrepo.RestaurantCatalogDTO, error) {
	if r.restaurantDB == nil {
		return nil, pkgErrors.ErrRepositoryNotInitialized
	}

	rows, err := r.restaurantDB.QueryContext(ctx, "SELECT id, name, address, status FROM RESTAURANTS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	restaurants := []customerrepo.RestaurantCatalogDTO{}
	for rows.Next() {
		var restaurant customerrepo.RestaurantCatalogDTO
		if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.Status); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *catalogRepo) ListRestaurantMenu(ctx context.Context, restaurantID uuid.UUID) ([]customerrepo.RestaurantMenuCatalogDTO, error) {
	if r.restaurantDB == nil {
		return nil, pkgErrors.ErrRepositoryNotInitialized
	}

	rows, err := r.restaurantDB.QueryContext(
		ctx,
		`SELECT order_item_id, restaurant_id, name, price, quantity, description
		 FROM RESTAURANT_MENU_ITEMS
		 WHERE restaurant_id = $1`,
		restaurantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menu := []customerrepo.RestaurantMenuCatalogDTO{}
	for rows.Next() {
		var item customerrepo.RestaurantMenuCatalogDTO
		if err := rows.Scan(&item.ID, &item.RestaurantID, &item.Name, &item.Price, &item.Quantity, &item.Description); err != nil {
			return nil, err
		}
		menu = append(menu, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menu, nil
}
