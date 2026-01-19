package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"customer/models"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	CourierID  uuid.UUID `json:"courier_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
}

type Filter struct {
	CustomerID *uuid.UUID
	CourierID  *uuid.UUID
	Status     string
}

type OrderItemInput struct {
	RestaurantItemID uuid.UUID
	Price            float64
	Quantity         int
}

type AcceptInput struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	CourierID  uuid.UUID
	Items      []OrderItemInput
}

type AcceptResult struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}

// update status (id, new status)

type Repository interface {
	Create(ctx context.Context, order Order) (Order, error)
	CreateWithItems(ctx context.Context, order Order, items []OrderItemInput) (Order, error)
	List(ctx context.Context, filter Filter) ([]Order, error)
	Accept(ctx context.Context, input AcceptInput) (AcceptResult, error)
}

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCourierNotFound  = errors.New("courier not found")
)

type postgresRepository struct {
	ordersDB    *sql.DB
	customersDB *sql.DB
	couriersDB  *sql.DB
}

func NewPostgresRepository(ordersDB, customersDB, couriersDB *sql.DB) Repository {
	return &postgresRepository{ordersDB: ordersDB, customersDB: customersDB, couriersDB: couriersDB}
}

func (r *postgresRepository) Create(ctx context.Context, order Order) (Order, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return Order{}, errors.New("orders repository not fully initialized")
	}

	now := time.Now().UTC()
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = now
	order.UpdatedAt = now
	if strings.TrimSpace(order.Status) == "" {
		order.Status = "created"
	}

	if _, err := r.ensureExists(ctx, r.customersDB, "SELECT 1 FROM customers WHERE emp_id = $1", order.CustomerID); err != nil {
		return Order{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", order.CourierID); err != nil {
		return Order{}, err
	}

	const insertQuery = `
        INSERT INTO ORDERS (emp_id, customer_id, courier_id, created_at, updated_at, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.ordersDB.ExecContext(ctx, insertQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) CreateWithItems(ctx context.Context, order Order, items []OrderItemInput) (Order, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return Order{}, errors.New("orders repository not fully initialized")
	}
	if len(items) == 0 {
		return Order{}, errors.New("items must not be empty")
	}

	now := time.Now().UTC()
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = now
	order.UpdatedAt = now
	if strings.TrimSpace(order.Status) == "" {
		order.Status = "created"
	}

	if _, err := r.ensureExists(ctx, r.customersDB, "SELECT 1 FROM customers WHERE emp_id = $1", order.CustomerID); err != nil {
		return Order{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", order.CourierID); err != nil {
		return Order{}, err
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertOrderQuery = `
		INSERT INTO ORDERS (emp_id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err = tx.ExecContext(ctx, insertOrderQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status); err != nil {
		return Order{}, err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (emp_id, order_id, restaurant_item_id, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range items {
		if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), order.ID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
			return Order{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) List(ctx context.Context, filter Filter) ([]Order, error) {
	query := `SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS`
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

	var result []Order
	for rows.Next() {
		var order Order
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

func (r *postgresRepository) Accept(ctx context.Context, input AcceptInput) (AcceptResult, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return AcceptResult{}, errors.New("orders repository not fully initialized")
	}
	if input.OrderID == uuid.Nil {
		return AcceptResult{}, errors.New("order_id must be a valid UUID")
	}
	if input.CustomerID == uuid.Nil {
		return AcceptResult{}, errors.New("customer_id must be a valid UUID")
	}
	if input.CourierID == uuid.Nil {
		return AcceptResult{}, errors.New("courier_id must be a valid UUID")
	}

	if _, err := r.ensureExists(ctx, r.customersDB, "SELECT 1 FROM customers WHERE emp_id = $1", input.CustomerID); err != nil {
		return AcceptResult{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", input.CourierID); err != nil {
		return AcceptResult{}, err
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return AcceptResult{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingStatus string
	statusQuery := "SELECT status FROM ORDERS WHERE emp_id = $1"
	scanErr := tx.QueryRowContext(ctx, statusQuery, input.OrderID).Scan(&existingStatus)
	if scanErr == nil && strings.EqualFold(existingStatus, string(models.OrderStatusKitchenAccepted)) {
		if err = tx.Commit(); err != nil {
			return AcceptResult{}, err
		}
		return AcceptResult{OrderID: input.OrderID, Status: existingStatus}, nil
	}
	if scanErr != nil && !errors.Is(scanErr, sql.ErrNoRows) {
		err = scanErr
		return AcceptResult{}, err
	}

	now := time.Now().UTC()
	const insertOrderQuery = `
		INSERT INTO ORDERS (emp_id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (emp_id) DO UPDATE
		SET status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`
	if _, err = tx.ExecContext(ctx, insertOrderQuery, input.OrderID, input.CustomerID, input.CourierID, now, now, string(models.OrderStatusKitchenAccepted)); err != nil {
		return AcceptResult{}, err
	}

	var itemsCount int
	countQuery := "SELECT COUNT(1) FROM ORDERS_ITEMS WHERE order_id = $1"
	if err = tx.QueryRowContext(ctx, countQuery, input.OrderID).Scan(&itemsCount); err != nil {
		return AcceptResult{}, err
	}
	if itemsCount == 0 && len(input.Items) > 0 {
		const insertItemQuery = `
			INSERT INTO ORDERS_ITEMS (emp_id, order_id, restaurant_item_id, price, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`
		for _, item := range input.Items {
			if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), input.OrderID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
				return AcceptResult{}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return AcceptResult{}, err
	}

	return AcceptResult{OrderID: input.OrderID, Status: string(models.OrderStatusKitchenAccepted)}, nil
}

func (r *postgresRepository) ensureExists(ctx context.Context, db *sql.DB, query string, id uuid.UUID) (bool, error) {
	var dummy int
	err := db.QueryRowContext(ctx, query, id).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		if strings.Contains(strings.ToLower(query), "customers") {
			return false, ErrCustomerNotFound
		}
		return false, ErrCourierNotFound
	}
	return false, err
}
