package orders

import (
	"context"
	"database/sql"
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

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, order Order) (Order, error) {
	now := time.Now().UTC()
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = now
	order.UpdatedAt = now
	if strings.TrimSpace(order.Status) == "" {
		order.Status = "created"
	}

	const insertQuery = `
        INSERT INTO ORDERS (empId, customer_id, courrier_id, created_at, updated_at, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.db.ExecContext(ctx, insertQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) List(ctx context.Context, filter Filter) ([]Order, error) {
	query := `SELECT empId, customer_id, courrier_id, created_at, updated_at, status FROM ORDERS`
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
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
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
