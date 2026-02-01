package usecase

import (
	"context"
	"strings"
	"time"

	"restaurant/internal/service"

	"github.com/Kabanya/YAFDS/pkg/models"

	"github.com/google/uuid"
)

type OrdersUseCase interface {
	ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error)
}

type ordersUseCase struct {
	// service service.OrdersService
	orderService service.OrdersService
}

// func NewOrdersUseCase(service service.OrdersService) OrdersUseCase {
// 	return &ordersUseCase{service: service}
// }

// func (u *ordersUseCase) ListOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]models.Order, error) {
// 	return u.service.ListOrdersByRestaurantID(ctx, restaurantID, status)
// }

func (u *ordersUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID) error {
	// if r.ordersDB == nil || r.customersDB == nil || r.couriersDB == nil {
	// 	return repositoryModels.AcceptResult{}, errors.New("orders repository not fully initialized")
	// }
	// if input.OrderID == uuid.Nil {
	// 	return repositoryModels.AcceptResult{}, errors.New("order_id must be a valid UUID")
	// }
	// if input.CustomerID == uuid.Nil {
	// 	return repositoryModels.AcceptResult{}, errors.New("customer_id must be a valid UUID")
	// }
	// if input.CourierID == uuid.Nil {
	// 	return repositoryModels.AcceptResult{}, errors.New("courier_id must be a valid UUID")
	// }

	// if err := r.ensureCustomerExists(ctx, input.CustomerID); err != nil {
	// 	return repositoryModels.AcceptResult{}, err
	// }
	// if err := r.ensureCourierExists(ctx, input.CourierID); err != nil {
	// 	return repositoryModels.AcceptResult{}, err
	// }

	// tx, err := r.ordersDB.BeginTx(ctx, nil)
	// if err != nil {
	// 	return repositoryModels.AcceptResult{}, err
	// }

	// defer func() { // выполнится после фукнции в горутине где мы живем. выполнится даже если закончится с exception
	// 	if err != nil { // если что-то пошло но так => rollback
	// 		_ = tx.Rollback()
	// 	}
	// }()

	var existingStatus string
	// statusQuery := "SELECT status FROM ORDERS WHERE id = $1"
	scanErr := tx.QueryRowContext(ctx, statusQuery, input.OrderID).Scan(&existingStatus)

	if scanErr != nil { //&& errors.Is(scanErr, sql.ErrNoRows) {
		err = scanErr
		return repositoryModels.AcceptResult{}, err
	}
	if strings.EqualFold(existingStatus, string(models.OrderStatusCustomerPaid)) {
		err != nil {
			return repositoryModels.AcceptResult{}, err
		}
		return repositoryModels.AcceptResult{OrderID: input.OrderID, Status: existingStatus}, nil
	}

	now := time.Now().UTC()
	status := input.Status
	if status == "" {
		status = models.OrderStatusKitchenAccepted
	}

	if _, err = tx.ExecContext(ctx, insertOrderQuery, input.OrderID, input.CustomerID, input.CourierID, now, now, string(status)); err != nil {
		return repositoryModels.AcceptResult{}, err
	}

	// var itemsCount int
	// countQuery := "SELECT COUNT(1) FROM ORDERS_ITEMS WHERE order_id = $1"
	// if err = tx.QueryRowContext(ctx, countQuery, input.OrderID).Scan(&itemsCount); err != nil {
	// 	return repositoryModels.AcceptResult{}, err
	// }
	// if itemsCount == 0 && len(input.Items) > 0 {
	// 	const insertItemQuery = `
	// 		INSERT INTO ORDERS_ITEMS (id, order_id, restaurant_item_id, price, quantity)
	// 		VALUES ($1, $2, $3, $4, $5)
	// 	`
	// 	for _, item := range input.Items {
	// 		if _, err = tx.ExecContext(ctx, insertItemQuery, uuid.New(), input.OrderID, item.RestaurantItemID, item.Price, item.Quantity); err != nil {
	// 			return repositoryModels.AcceptResult{}, err
	// 		}
	// 	}
	// }

	if err = tx.Commit(); err != nil {
		return repositoryModels.AcceptResult{}, err
	}

	return repositoryModels.AcceptResult{OrderID: input.OrderID, Status: string(status)}, nil
}
