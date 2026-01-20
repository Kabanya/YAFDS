package repository

import (
	"customer/pkg/utils"
	"database/sql"
	"restaurant/pkg/models"

	"github.com/google/uuid"
)

type RestaurantMenuItemsRepo interface {
	ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error)
	UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error
}

type restaurantMenuItemsRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewRestaurantMenuItemsRepo(db *sql.DB) *restaurantMenuItemsRepo {
	return &restaurantMenuItemsRepo{db: db}
}

func (r *restaurantMenuItemsRepo) ShowMenuItemsByRestaurantID(restaurantID uuid.UUID) ([]models.MenuItem, error) {
	logger, err := utils.Logger()
	if err != nil {
		return nil, err
	}
	sqlStatement := `
	       SELECT order_item_id, restaurant_id, name, price, quantity, description
	       FROM restaurant_menu_items
	       WHERE restaurant_id = $1
       `
	rows, err := r.db.Query(sqlStatement, restaurantID)
	if err != nil {
		logger.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var menuItems []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		err := rows.Scan(&item.OrderItemID, &item.RestaurantID, &item.Name, &item.Price, &item.Quantity, &item.Description)
		if err != nil {
			logger.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		menuItems = append(menuItems, item)
	}
	if err = rows.Err(); err != nil {
		logger.Printf("Rows error: %v", err)
		return nil, err
	}
	return menuItems, nil
}

func (r *restaurantMenuItemsRepo) UploadMenuItemsByRestaurantID(menuItem models.MenuItem) error {
	logger, err := utils.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
	       INSERT INTO restaurant_menu_items (order_item_id, restaurant_id, name, price, quantity, description)
	       VALUES ($1, $2, $3, $4, $5, $6)
	   `

	_, err = r.db.Exec(sqlStatement, menuItem.OrderItemID, menuItem.RestaurantID, menuItem.Name, menuItem.Price, menuItem.Quantity, menuItem.Description)
	if err != nil {
		logger.Printf("Failed to execute insert: %v", err)
		return err
	}

	logger.Printf("Successfully inserted menu item: %s for restaurant: %s", menuItem.Name, menuItem.RestaurantID)
	return nil
}
