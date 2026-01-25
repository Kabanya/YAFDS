package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"

	"github.com/google/uuid"
)

type postgresRepository struct {
	ordersDB    *sql.DB
	customersDB *sql.DB
	couriersDB  *sql.DB
}

func NewPostgresRepository(ordersDB, customersDB, couriersDB *sql.DB) repositoryModels.Order {
	return &postgresRepository{ordersDB: ordersDB, customersDB: customersDB, couriersDB: couriersDB}
}

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCourierNotFound  = errors.New("courier not found")
	ErrOrderNotFound    = errors.New("order not found")
)

func (r *postgresRepository) Create(ctx context.Context, order models.Order) (models.Order, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return models.Order{}, errors.New("orders repository not fully initialized")
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
		return models.Order{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", order.CourierID); err != nil {
		return models.Order{}, err
	}

	const insertQuery = `
        INSERT INTO ORDERS (emp_id, customer_id, courier_id, created_at, updated_at, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.ordersDB.ExecContext(ctx, insertQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) CreateWithItems(ctx context.Context, order models.Order, items []repositoryModels.OrderItemInput) (models.Order, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return models.Order{}, errors.New("orders repository not fully initialized")
	}
	if len(items) == 0 {
		return models.Order{}, errors.New("items must not be empty")
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
		return models.Order{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", order.CourierID); err != nil {
		return models.Order{}, err
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return models.Order{}, err
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
		return models.Order{}, err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (emp_id, order_id, restaurant_item_id, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range items {
		if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), order.ID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
			return models.Order{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) List(ctx context.Context, filter repositoryModels.Filter) ([]models.Order, error) {
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

func (r *postgresRepository) Get(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	if r.ordersDB == nil {
		return models.Order{}, errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return models.Order{}, errors.New("order_id must be a valid UUID")
	}

	var order models.Order
	query := `SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS WHERE emp_id = $1`
	err := r.ordersDB.QueryRowContext(ctx, query, orderID).Scan(&order.ID, &order.CustomerID, &order.CourierID, &order.CreatedAt, &order.UpdatedAt, &order.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, ErrOrderNotFound
		}
		return models.Order{}, err
	}
	return order, nil
}

func (r *postgresRepository) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	if r.ordersDB == nil {
		return "", errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return "", errors.New("order_id must be a valid UUID")
	}

	var status string
	query := "SELECT status FROM ORDERS WHERE emp_id = $1"
	if err := r.ordersDB.QueryRowContext(ctx, query, orderID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrOrderNotFound
		}
		return "", err
	}
	return models.OrderStatus(status), nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	if r.ordersDB == nil {
		return errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return errors.New("order_id must be a valid UUID")
	}

	res, err := r.ordersDB.ExecContext(ctx, "UPDATE ORDERS SET status = $1, updated_at = $2 WHERE emp_id = $3", string(status), time.Now().UTC(), orderID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return ErrOrderNotFound
	}
	return err
}

func (r *postgresRepository) GetOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	if r.ordersDB == nil {
		return 0, errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return 0, errors.New("order_id must be a valid UUID")
	}

	var total sql.NullFloat64
	query := "SELECT SUM(price * quantity) FROM ORDERS_ITEMS WHERE order_id = $1"
	if err := r.ordersDB.QueryRowContext(ctx, query, orderID).Scan(&total); err != nil {
		return 0, err
	}
	if !total.Valid {
		return 0, nil
	}
	return total.Float64, nil
}

func (r *postgresRepository) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	if r.customersDB == nil {
		return "", errors.New("customers repository not fully initialized")
	}
	if customerID == uuid.Nil {
		return "", errors.New("customer_id must be a valid UUID")
	}

	var wallet string
	query := "SELECT wallet_address FROM CUSTOMERS WHERE emp_id = $1"
	if err := r.customersDB.QueryRowContext(ctx, query, customerID).Scan(&wallet); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCustomerNotFound
		}
		return "", err
	}
	if strings.TrimSpace(wallet) == "" {
		return "", errors.New("wallet_address is empty")
	}
	return wallet, nil
}

func (r *postgresRepository) Accept(ctx context.Context, input repositoryModels.AcceptInput) (repositoryModels.AcceptResult, error) {
	if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
		return repositoryModels.AcceptResult{}, errors.New("orders repository not fully initialized")
	}
	if input.OrderID == uuid.Nil {
		return repositoryModels.AcceptResult{}, errors.New("order_id must be a valid UUID")
	}
	if input.CustomerID == uuid.Nil {
		return repositoryModels.AcceptResult{}, errors.New("customer_id must be a valid UUID")
	}
	if input.CourierID == uuid.Nil {
		return repositoryModels.AcceptResult{}, errors.New("courier_id must be a valid UUID")
	}

	if _, err := r.ensureExists(ctx, r.customersDB, "SELECT 1 FROM customers WHERE emp_id = $1", input.CustomerID); err != nil {
		return repositoryModels.AcceptResult{}, err
	}
	if _, err := r.ensureExists(ctx, r.couriersDB, "SELECT 1 FROM couriers WHERE emp_id = $1", input.CourierID); err != nil {
		return repositoryModels.AcceptResult{}, err
	}

	tx, err := r.ordersDB.BeginTx(ctx, nil)
	if err != nil {
		return repositoryModels.AcceptResult{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingStatus string
	statusQuery := "SELECT status FROM ORDERS WHERE emp_id = $1"
	scanErr := tx.QueryRowContext(ctx, statusQuery, input.OrderID).Scan(&existingStatus)
	if scanErr == nil && (strings.EqualFold(existingStatus, string(models.OrderStatusKitchenAccepted)) || strings.EqualFold(existingStatus, string(models.OrderStatusKitchenDenied))) {
		if err = tx.Commit(); err != nil {
			return repositoryModels.AcceptResult{}, err
		}
		return repositoryModels.AcceptResult{OrderID: input.OrderID, Status: existingStatus}, nil
	}
	if scanErr != nil && !errors.Is(scanErr, sql.ErrNoRows) {
		err = scanErr
		return repositoryModels.AcceptResult{}, err
	}

	now := time.Now().UTC()
	status := input.Status
	if status == "" {
		status = models.OrderStatusKitchenAccepted
	}
	const insertOrderQuery = `
		INSERT INTO ORDERS (emp_id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (emp_id) DO UPDATE
		SET status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`
	if _, err = tx.ExecContext(ctx, insertOrderQuery, input.OrderID, input.CustomerID, input.CourierID, now, now, string(status)); err != nil {
		return repositoryModels.AcceptResult{}, err
	}

	var itemsCount int
	countQuery := "SELECT COUNT(1) FROM ORDERS_ITEMS WHERE order_id = $1"
	if err = tx.QueryRowContext(ctx, countQuery, input.OrderID).Scan(&itemsCount); err != nil {
		return repositoryModels.AcceptResult{}, err
	}
	if itemsCount == 0 && len(input.Items) > 0 {
		const insertItemQuery = `
			INSERT INTO ORDERS_ITEMS (emp_id, order_id, restaurant_item_id, price, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`
		for _, item := range input.Items {
			if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), input.OrderID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
				return repositoryModels.AcceptResult{}, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return repositoryModels.AcceptResult{}, err
	}

	return repositoryModels.AcceptResult{OrderID: input.OrderID, Status: string(status)}, nil
}

func (r *postgresRepository) AddItem(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
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
	if err = tx.QueryRowContext(ctx, "SELECT 1 FROM ORDERS WHERE emp_id = $1", orderID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderNotFound
		}
		return err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (emp_id, order_id, restaurant_item_id, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), orderID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, "UPDATE ORDERS SET updated_at = $1 WHERE emp_id = $2", time.Now().UTC(), orderID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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
