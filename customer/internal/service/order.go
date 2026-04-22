package service

import (
	"context"
	"errors"
	"strings"
	"time"

	customerrepo "github.com/Kabanya/YAFDS/customer/internal/repository"
	"github.com/Kabanya/YAFDS/pkg/order/domain"
	ordermodels "github.com/Kabanya/YAFDS/pkg/order/usecase"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (domain.Order, error)
	CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []ordermodels.OrderItemInput) (domain.Order, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
	PayOrder(ctx context.Context, orderID uuid.UUID) error
}

type orderService struct {
	repo customerrepo.CustomerOrderRepository
}

func NewOrderService(repo customerrepo.CustomerOrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (domain.Order, error) {
	order := newCustomerOrderDTO(customerID, courierID)
	created, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}
	return orderDTOToDomain(created), nil
}

func (s *orderService) CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []ordermodels.OrderItemInput) (domain.Order, error) {
	if len(items) == 0 {
		return domain.Order{}, errors.New("items must not be empty")
	}

	order := newCustomerOrderDTO(customerID, courierID)
	created, err := s.repo.CreateOrderWithItems(ctx, order, orderItemsToDTO(items))
	if err != nil {
		return domain.Order{}, err
	}
	return orderDTOToDomain(created), nil
}

func (s *orderService) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	walletAddress, err := s.repo.GetCustomerWalletAddress(ctx, customerID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(walletAddress) == "" {
		return "", errors.New("wallet_address is empty")
	}
	return walletAddress, nil
}

func (s *orderService) PayOrder(ctx context.Context, orderID uuid.UUID) error {
	return s.repo.UpdateOrderStatus(ctx, orderID, string(domain.OrderStatusCustomerPaid))
}

func newCustomerOrderDTO(customerID uuid.UUID, courierID uuid.UUID) customerrepo.OrderDTO {
	now := time.Now().UTC()
	return customerrepo.OrderDTO{
		ID:         uuid.New(),
		CustomerID: customerID,
		CourierID:  courierID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Status:     string(domain.OrderStatusCustomerCreated),
	}
}

func orderDTOToDomain(order customerrepo.OrderDTO) domain.Order {
	return domain.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		CourierID:  order.CourierID,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Status:     domain.OrderStatus(order.Status),
	}
}

func orderItemsToDTO(items []ordermodels.OrderItemInput) []customerrepo.OrderItemDTO {
	result := make([]customerrepo.OrderItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, customerrepo.OrderItemDTO{
			RestaurantItemID: item.RestaurantItemID,
			Price:            item.Price,
			Quantity:         item.Quantity,
		})
	}
	return result
}
