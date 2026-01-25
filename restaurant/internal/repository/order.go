package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrdersRepo interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
}

type ordersRepo struct {
	ordersDB     *sql.DB
	restaurantDB *sql.DB
}

func NewOrdersRepo(ordersDB, restaurantDB *sql.DB) OrdersRepo {
	return &ordersRepo{ordersDB: ordersDB, restaurantDB: restaurantDB}
}

func (r *ordersRepo) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
	if r.ordersDB == nil || r.restaurantDB == nil {
		return nil, errors.New("orders repository not fully initialized")
	}

	rows, err := r.restaurantDB.QueryContext(ctx, `
		SELECT order_item_id
		FROM restaurant_menu_items
		WHERE restaurant_id = $1
	`, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(itemIDs) == 0 {
		return []models.Order{}, nil
	}

	placeholders := make([]string, len(itemIDs))
	args := make([]any, 0, len(itemIDs)+1)
	for i, id := range itemIDs {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args = append(args, id)
	}

	query := `
		SELECT DISTINCT o.emp_id, o.customer_id, o.courier_id, o.created_at, o.updated_at, o.status
		FROM ORDERS o
		JOIN ORDERS_ITEMS oi ON oi.order_id = o.emp_id
		WHERE oi.restaurant_item_id IN (` + strings.Join(placeholders, ",") + `)
	`

	if status != "" {
		query += " AND o.status = $" + strconv.Itoa(len(args)+1)
		args = append(args, status)
	}

	query += " ORDER BY o.created_at DESC"

	orderRows, err := r.ordersDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer orderRows.Close()

	var result []models.Order
	for orderRows.Next() {
		var order models.Order
		if err := orderRows.Scan(&order.ID, &order.CustomerID, &order.CourierID, &order.CreatedAt, &order.UpdatedAt, &order.Status); err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	if err := orderRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
