package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/google/uuid"

	errorsPkg "github.com/Kabanya/YAFDS/pkg/common/errors"
	repoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type OrderRepo interface { //   dto - data transfer object. dto похож на model но другое.
	//                          Если совпадают, то приоритет модели по возможности не плодим dto
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)            // PKG (all)
	ListOrders(ctx context.Context, filter repoModels.Filter) ([]models.Order, error) // PKG (all)

	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)         // PKG (all)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error // PKG (all)
	// ЕСЛИ Прям нужна целостность данных (транзакции, которые можно только на уровне бд открывать) то пишем в repository слой, а так не выебываемся

	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]models.MenuItem, error)               // PKG (customer, restaurant)
	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repoModels.OrderItemInput) error // PKG (customer, restaurant)
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error  // PKG (customer, restaurant)
}

type postgresRepository struct {
	ordersDB    *sql.DB
	customersDB *sql.DB
	couriersDB  *sql.DB
}

func NewPostgresRepository(ordersDB, customersDB, couriersDB *sql.DB) OrderRepo {
	return &postgresRepository{ordersDB: ordersDB, customersDB: customersDB, couriersDB: couriersDB}
}

func (r *postgresRepository) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	if r.ordersDB == nil {
		return models.Order{}, errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return models.Order{}, errors.New("order_id must be a valid UUID")
	}

	var order models.Order
	query := `SELECT id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS WHERE id = $1`
	err := r.ordersDB.QueryRowContext(ctx, query, orderID).Scan(&order.ID, &order.CustomerID, &order.CourierID, &order.CreatedAt, &order.UpdatedAt, &order.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, errorsPkg.ErrOrderNotFound
		}
		return models.Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) ListOrders(ctx context.Context, filter repoModels.Filter) ([]models.Order, error) {
	query := `SELECT id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS`
	var args []any
	var where []string

	if filter.CustomerID != nil {
		where = append(where, "customer_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, *filter.CustomerID)
	}
	if filter.CourierID != nil {
		where = append(where, "courrier_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, *filter.CourierID)
	}
	if filter.Status != "" {
		where = append(where, "status = $"+strconv.Itoa(len(args)+1))
		args = append(args, filter.Status)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY created_at DESC "

	rows, err := r.ordersDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.CourierID, &order.CreatedAt, &order.UpdatedAt, &order.Status); err != nil {
			return nil, err
		}
		result = append(result, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// func (r *postgresRepository) AcceptOrder(ctx context.Context, input repoModels.AcceptInput) (repoModels.AcceptResult, error) {
// 	return repoModels.AcceptResult{}, nil
// }

func (r *postgresRepository) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	if r.ordersDB == nil {
		return "", errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return "", errors.New("order_id must be a valid UUID")
	}

	var status string
	query := "SELECT status FROM ORDERS WHERE id = $1"
	if err := r.ordersDB.QueryRowContext(ctx, query, orderID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errorsPkg.ErrOrderNotFound
		}
		return "", err
	}
	return models.OrderStatus(status), nil
}

func (r *postgresRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	const query = "UPDATE ORDERS SET status = $1, updated_at = NOW() WHERE id = $2"
	res, err := r.ordersDB.ExecContext(ctx, query, string(status), orderID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 { // если не обновилось ни одной строки => заказа нет
		return errorsPkg.ErrOrderNotFound
	}
	return err
}

func (r *postgresRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]models.MenuItem, error) {
	if r.ordersDB == nil {
		return nil, errors.New("orders repository not fully initialized")
	}

	query := `
		SELECT restaurant_item_id, price, quantity
		FROM ORDERS_ITEMS
		WHERE order_id = $1
	`
	rows, err := r.ordersDB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.RestaurantID, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *postgresRepository) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repoModels.OrderItemInput) error {
	if r.ordersDB == nil {
		return errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return errors.New("order_id must be a valid UUID")
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var exists int
	if err = tx.QueryRowContext(ctx, "SELECT 1 FROM ORDERS WHERE id = $1", orderID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorsPkg.ErrOrderNotFound
		}
		return err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (id, order_id, restaurant_item_id, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), orderID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, "UPDATE ORDERS SET updated_at = $1 WHERE id = $2", time.Now().UTC(), orderID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	if r.ordersDB == nil {
		return errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return errors.New("order_id must be a valid UUID")
	}
	if restaurantItemID == uuid.Nil {
		return errors.New("restaurant_item_id must be a valid UUID")
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var exists int
	if err = tx.QueryRowContext(ctx, "SELECT 1 FROM ORDERS WHERE id = $1", orderID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorsPkg.ErrOrderNotFound
		}
		return err
	}

	res, err := tx.ExecContext(ctx, "DELETE FROM ORDERS_ITEMS WHERE order_id = $1 AND restaurant_item_id = $2", orderID, restaurantItemID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		err = errors.New("item not found in order")
		return err
	}

	if _, err = tx.ExecContext(ctx, "UPDATE ORDERS SET updated_at = $1 WHERE id = $2", time.Now().UTC(), orderID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
