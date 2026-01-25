package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"

	"github.com/google/uuid"
)

// НА УРОВНЕ USECASE ходим за продуктовыми логикой
var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrWalletUnavailable       = errors.New("wallet service unavailable")
	ErrInsufficientFunds       = errors.New("insufficient funds")
)

type WalletClient interface {
	CheckAndDebit(ctx context.Context, walletAddress string, amount float64) (bool, error)
}

type OrderUseCase interface {
	Pay(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) (models.OrderStatus, error)
	ChangeStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus) (models.OrderStatus, error)
}

type orderUseCase struct {
	repo   repositoryModels.Order
	wallet WalletClient
}

func NewOrderUseCase(repo repositoryModels.Order, wallet WalletClient) OrderUseCase {
	return &orderUseCase{repo: repo, wallet: wallet}
}

func (u *orderUseCase) Pay(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) (models.OrderStatus, error) {
	if orderID == uuid.Nil {
		return "", errors.New("order_id must be UUID")
	}
	if customerID == uuid.Nil {
		return "", errors.New("customer_id must be UUID")
	}
	if u.wallet == nil {
		return "", ErrWalletUnavailable
	}

	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return "", err
	}
	if order.CustomerID != customerID {
		return "", errors.New("order does not belong to customer")
	}

	current := normalizeStatus(models.OrderStatus(order.Status))
	if current != models.OrderStatusCustomerCreated {
		return current, fmt.Errorf("order status must be %s to pay", models.OrderStatusCustomerCreated)
	}

	walletAddress, err := u.repo.GetCustomerWalletAddress(ctx, customerID)
	if err != nil {
		return "", err
	}
	amount, err := u.repo.GetOrderTotal(ctx, orderID)
	if err != nil {
		return "", err
	}
	if amount <= 0 {
		return "", errors.New("order total must be positive")
	}

	paid, err := u.wallet.CheckAndDebit(ctx, walletAddress, amount)
	if err != nil {
		return "", err
	}
	if !paid {
		if _, err := u.ChangeStatus(ctx, orderID, models.OrderStatusCustomerCancelled); err != nil {
			return "", err
		}
		return models.OrderStatusCustomerCancelled, ErrInsufficientFunds
	}

	return u.ChangeStatus(ctx, orderID, models.OrderStatusCustomerPaid)
}

func (u *orderUseCase) ChangeStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus) (models.OrderStatus, error) {
	if orderID == uuid.Nil {
		return "", errors.New("order_id must be UUID")
	}
	current, err := u.repo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return "", err
	}
	current = normalizeStatus(current)
	if current == newStatus {
		return current, nil
	}
	if !isTransitionAllowed(current, newStatus) {
		return current, fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, current, newStatus)
	}
	if err := u.repo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		return "", err
	}
	return newStatus, nil
}

func isTransitionAllowed(from models.OrderStatus, to models.OrderStatus) bool {
	allowed := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusCustomerCreated:    {models.OrderStatusCustomerPaid, models.OrderStatusCustomerCancelled},
		models.OrderStatusCustomerPaid:       {models.OrderStatusKitchenAccepted, models.OrderStatusKitchenDenied},
		models.OrderStatusKitchenAccepted:    {models.OrderStatusKitchenPreparing},
		models.OrderStatusKitchenDenied:      {models.OrderStatusCourierRefunded},
		models.OrderStatusKitchenPreparing:   {models.OrderStatusDeliveryPending},
		models.OrderStatusDeliveryPending:    {models.OrderStatusDeliveryPicking, models.OrderStatusDeliveryDenied},
		models.OrderStatusDeliveryPicking:    {models.OrderStatusDeliveryDelivering, models.OrderStatusDeliveryDenied},
		models.OrderStatusDeliveryDelivering: {models.OrderStatusOrderCompleted, models.OrderStatusDeliveryRefunded},
		models.OrderStatusDeliveryDenied:     {models.OrderStatusDeliveryRefunded},
		models.OrderStatusCourierRefunded:    {models.OrderStatusOrderCompleted},
		models.OrderStatusDeliveryRefunded:   {models.OrderStatusOrderCompleted},
	}

	next, ok := allowed[from]
	if !ok {
		return false
	}
	for _, candidate := range next {
		if candidate == to {
			return true
		}
	}
	return false
}

func normalizeStatus(status models.OrderStatus) models.OrderStatus {
	normalized := strings.ToUpper(strings.TrimSpace(string(status)))
	if mapped, ok := statusNormalizationMap[normalized]; ok {
		return mapped
	}
	return status
}

var statusNormalizationMap = map[string]models.OrderStatus{
	"CREATED": models.OrderStatusCustomerCreated,
	string(models.OrderStatusCustomerCreated):    models.OrderStatusCustomerCreated,
	string(models.OrderStatusCustomerPaid):       models.OrderStatusCustomerPaid,
	string(models.OrderStatusCustomerCancelled):  models.OrderStatusCustomerCancelled,
	string(models.OrderStatusKitchenAccepted):    models.OrderStatusKitchenAccepted,
	string(models.OrderStatusKitchenDenied):      models.OrderStatusKitchenDenied,
	string(models.OrderStatusKitchenPreparing):   models.OrderStatusKitchenPreparing,
	string(models.OrderStatusCourierRefunded):    models.OrderStatusCourierRefunded,
	string(models.OrderStatusDeliveryPending):    models.OrderStatusDeliveryPending,
	string(models.OrderStatusDeliveryPicking):    models.OrderStatusDeliveryPicking,
	string(models.OrderStatusDeliveryDenied):     models.OrderStatusDeliveryDenied,
	string(models.OrderStatusDeliveryRefunded):   models.OrderStatusDeliveryRefunded,
	string(models.OrderStatusDeliveryDelivering): models.OrderStatusDeliveryDelivering,
	string(models.OrderStatusOrderCompleted):     models.OrderStatusOrderCompleted,
}
