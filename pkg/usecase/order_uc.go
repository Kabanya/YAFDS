package usecase

import (
	"context"
	"errors"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	service "github.com/Kabanya/YAFDS/pkg/service"

	"github.com/google/uuid"
)

// НА УРОВНЕ USECASE ходим за продуктовыми логикой
var (
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrWalletUnavailable       = errors.New("wallet service unavailable")
	ErrInsufficientFunds       = errors.New("insufficient funds")
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (models.Order, error)
	CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput) (models.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error)
	AcceptOrder(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput, status models.OrderStatus) (repositoryModels.AcceptResult, error)
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error
	RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error
}

type orderUseCase struct {
	serviceOrder service.OrderService
}

func NewOrderUseCase(serviceOrder service.OrderService) OrderUseCase {
	return &orderUseCase{serviceOrder: serviceOrder}
}

// type UserUseCase interface {
// 	Register(uuid.UUID, string, string, string, string) error
// 	Login(walletAddress string, password string) (models.LoginResponse, error)
// }

// type userUseCase struct {
// 	service service.UserService
// }

// func NewUserUseCase(service service.UserService) UserUseCase {
// 	return &userUseCase{service: service}
// }

// func (u *userUseCase) Register(id uuid.UUID, name string, walletAddress string, address string, password string) error {
// 	return u.service.Register(id, name, walletAddress, address, password)
// }

// func (u *userUseCase) Login(walletAddress string, password string) (models.LoginResponse, error) {
// 	return u.service.Login(walletAddress, password)
// }
