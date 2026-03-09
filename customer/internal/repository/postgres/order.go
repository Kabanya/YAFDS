package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	pkgErrors "github.com/Kabanya/YAFDS/pkg/common/errors"
	"github.com/Kabanya/YAFDS/pkg/models"
	repoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	pkgRepo "github.com/Kabanya/YAFDS/pkg/repository/postgres"

	"github.com/google/uuid"
)

type OrderUserRepo interface {
	CreateOrder(ctx context.Context, order repoModels.Filter) (models.Order, error)                                        // Customer
	CreateOrderWithItems(ctx context.Context, order models.Order, items []repoModels.OrderItemInput) (models.Order, error) // Customer
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)                                    // Customer
	PayOrder(ctx context.Context, orderID uuid.UUID) error                                                                 // Customer
}

type postgresUserRepo struct {
	customerDB *sql.DB
	courierDB  *sql.DB
	orderDB    *sql.DB
}

func NewPostgresUserRepo(customerDB, courierDB, orderDB *sql.DB) OrderUserRepo {
	return &postgresUserRepo{customerDB: customerDB, courierDB: courierDB, orderDB: orderDB}
}

func (r *postgresUserRepo) CreateOrder(ctx context.Context, filter repoModels.Filter) (models.Order, error) {
	if r.orderDB == nil || r.customerDB == nil || r.courierDB == nil {
		return models.Order{}, pkgErrors.ErrRepositoryNotInitialized
	}

	now := time.Now().UTC()
	order := models.Order{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Status:    models.OrderStatusCustomerCreated,
	}

	if filter.CustomerID != nil {
		order.CustomerID = *filter.CustomerID
	}
	if filter.CourierID != nil {
		order.CourierID = *filter.CourierID
	}
	if filter.Status != "" {
		order.Status = models.OrderStatus(filter.Status)
	}

	if err := pkgRepo.EnsureCustomerExists(ctx, r.customerDB, order.CustomerID); err != nil {
		return models.Order{}, err
	}
	if err := pkgRepo.EnsureCourierExists(ctx, r.courierDB, order.CourierID); err != nil {
		return models.Order{}, err
	}

	const insertQuery = `
		INSERT INTO ORDERS (id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.orderDB.ExecContext(ctx, insertQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, string(order.Status))

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *postgresUserRepo) CreateOrderWithItems(ctx context.Context, order models.Order, items []repoModels.OrderItemInput) (models.Order, error) {
	if r.orderDB == nil || r.customerDB == nil || r.courierDB == nil {
		return models.Order{}, pkgErrors.ErrRepositoryNotInitialized
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
	if strings.TrimSpace(string(order.Status)) == "" {
		order.Status = models.OrderStatusCustomerCreated
	}

	if err := pkgRepo.EnsureCustomerExists(ctx, r.customerDB, order.CustomerID); err != nil {
		return models.Order{}, err
	}
	if err := pkgRepo.EnsureCourierExists(ctx, r.courierDB, order.CourierID); err != nil {
		return models.Order{}, err
	}

	tx, err := r.orderDB.BeginTx(ctx, nil)
	if err != nil {
		return models.Order{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertOrderQuery = `
		INSERT INTO ORDERS (id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if _, err = tx.ExecContext(ctx, insertOrderQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, string(order.Status)); err != nil {
		return models.Order{}, err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (id, order_id, restaurant_item_id, price, quantity)
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

func (r *postgresUserRepo) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	if r.customerDB == nil {
		return "", errors.New("customers repository not fully initialized")
	}
	if customerID == uuid.Nil {
		return "", errors.New("customer_id must be a valid UUID")
	}

	var wallet string
	query := "SELECT wallet_address FROM CUSTOMERS WHERE id = $1"
	if err := r.customerDB.QueryRowContext(ctx, query, customerID).Scan(&wallet); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", pkgErrors.ErrCustomerNotFound
		}
		return "", err
	}
	if strings.TrimSpace(wallet) == "" {
		return "", errors.New("wallet_address is empty")
	}
	return wallet, nil
}

func (r *postgresUserRepo) PayOrder(ctx context.Context, orderID uuid.UUID) error {
	if r.orderDB == nil {
		return errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return errors.New("order_id must be a valid UUID")
	}

	var currentStatus string
	statusQuery := "SELECT status FROM ORDERS WHERE id = $1"
	err := r.orderDB.QueryRowContext(ctx, statusQuery, orderID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkgErrors.ErrOrderNotFound
		}
		return err
	}

	if currentStatus != string(models.OrderStatusCustomerCreated) {
		return errors.New("order can only be paid when in CUSTOMER_CREATED status")
	}

	updateQuery := "UPDATE ORDERS SET status = $1, updated_at = $2 WHERE id = $3"
	_, err = r.orderDB.ExecContext(ctx, updateQuery, string(models.OrderStatusCustomerPaid), time.Now().UTC(), orderID)
	if err != nil {
		return err
	}
	return nil
}
