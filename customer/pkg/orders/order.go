package orders

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

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

type Repository interface {
	Create(ctx context.Context, order Order) (Order, error)
	List(ctx context.Context, filter Filter) ([]Order, error)
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
