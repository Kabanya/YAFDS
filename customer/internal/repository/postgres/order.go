package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	courierexistence "github.com/Kabanya/YAFDS/courier/pkg/repository"
	customerrepo "github.com/Kabanya/YAFDS/customer/internal/repository"
	customerexistence "github.com/Kabanya/YAFDS/customer/pkg/repository"
	pkgErrors "github.com/Kabanya/YAFDS/pkg/errors"
	"github.com/google/uuid"
)

type OrderUserRepo interface {
	customerrepo.CustomerOrderRepository
}

type postgresUserRepo struct {
	customerDB *sql.DB
	courierDB  *sql.DB
	orderDB    *sql.DB
}

func NewPostgresUserRepo(customerDB, courierDB, orderDB *sql.DB) OrderUserRepo {
	return &postgresUserRepo{customerDB: customerDB, courierDB: courierDB, orderDB: orderDB}
}

func (r *postgresUserRepo) CreateOrder(ctx context.Context, order customerrepo.OrderDTO) (customerrepo.OrderDTO, error) {
	if r.orderDB == nil || r.customerDB == nil || r.courierDB == nil {
		return customerrepo.OrderDTO{}, pkgErrors.ErrRepositoryNotInitialized
	}

	if err := customerexistence.EnsureCustomerExists(ctx, r.customerDB, order.CustomerID); err != nil {
		return customerrepo.OrderDTO{}, err
	}
	if err := courierexistence.EnsureCourierExists(ctx, r.courierDB, order.CourierID); err != nil {
		return customerrepo.OrderDTO{}, err
	}

	const insertQuery = `
		INSERT INTO ORDERS (id, customer_id, courier_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.orderDB.ExecContext(ctx, insertQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status)

	if err != nil {
		return customerrepo.OrderDTO{}, err
	}

	return order, nil
}

func (r *postgresUserRepo) CreateOrderWithItems(ctx context.Context, order customerrepo.OrderDTO, items []customerrepo.OrderItemDTO) (customerrepo.OrderDTO, error) {
	if r.orderDB == nil || r.customerDB == nil || r.courierDB == nil {
		return customerrepo.OrderDTO{}, pkgErrors.ErrRepositoryNotInitialized
	}

	if err := customerexistence.EnsureCustomerExists(ctx, r.customerDB, order.CustomerID); err != nil {
		return customerrepo.OrderDTO{}, err
	}
	if err := courierexistence.EnsureCourierExists(ctx, r.courierDB, order.CourierID); err != nil {
		return customerrepo.OrderDTO{}, err
	}

	tx, err := r.orderDB.BeginTx(ctx, nil)
	if err != nil {
		return customerrepo.OrderDTO{}, err
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
	if _, err = tx.ExecContext(ctx, insertOrderQuery, order.ID, order.CustomerID, order.CourierID, order.CreatedAt, order.UpdatedAt, order.Status); err != nil {
		return customerrepo.OrderDTO{}, err
	}

	const insertItemQuery = `
		INSERT INTO ORDERS_ITEMS (id, order_id, restaurant_item_id, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range items {
		if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), order.ID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
			return customerrepo.OrderDTO{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return customerrepo.OrderDTO{}, err
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
	return wallet, nil
}

func (r *postgresUserRepo) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	if r.orderDB == nil {
		return errors.New("orders repository not fully initialized")
	}
	if orderID == uuid.Nil {
		return errors.New("order_id must be a valid UUID")
	}

	updateQuery := "UPDATE ORDERS SET status = $1, updated_at = $2 WHERE id = $3"
	res, err := r.orderDB.ExecContext(ctx, updateQuery, status, time.Now().UTC(), orderID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return pkgErrors.ErrOrderNotFound
	}
	return nil
}
