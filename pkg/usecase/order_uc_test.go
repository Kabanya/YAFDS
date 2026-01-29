package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/Kabanya/YAFDS/pkg/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, customerID string, courierID string, status models.OrderStatus) (service.CreateOrderResponse, error) {
	args := m.Called(ctx, customerID, courierID, status)
	return args.Get(0).(service.CreateOrderResponse), args.Error(1)
}

func (m *MockOrderService) CreateOrderWithItems(ctx context.Context, customerID string, courierID string, status models.OrderStatus, items []repositoryModels.OrderItemInput) (service.CreateOrderResponse, error) {
	args := m.Called(ctx, customerID, courierID, status, items)
	return args.Get(0).(service.CreateOrderResponse), args.Error(1)
}

func (m *MockOrderService) ListOrders(ctx context.Context, filter repositoryModels.Filter) ([]models.Order, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderService) AcceptOrder(ctx context.Context, orderID string, customerID string, courierID string, items []repositoryModels.OrderItemInput, status models.OrderStatus) (repositoryModels.AcceptResult, error) {
	args := m.Called(ctx, orderID, customerID, courierID, items, status)
	return args.Get(0).(repositoryModels.AcceptResult), args.Error(1)
}

func (m *MockOrderService) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.OrderStatus), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderService) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderService) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	args := m.Called(ctx, customerID)
	return args.String(0), args.Error(1)
}

func (m *MockOrderService) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
	args := m.Called(ctx, orderID, item)
	return args.Error(0)
}

func (m *MockOrderService) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	args := m.Called(ctx, orderID, restaurantItemID)
	return args.Error(0)
}

func createTestOrder() models.Order {
	return models.Order{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		CourierID:  uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     models.OrderStatusCustomerCreated,
	}
}

func createTestOrderItemInput() repositoryModels.OrderItemInput {
	return repositoryModels.OrderItemInput{
		RestaurantItemID: uuid.New(),
		Price:            10.5,
		Quantity:         2,
	}
}

func TestOrderUseCase_CreateOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()

	createdOrder := createTestOrder()
	response := service.CreateOrderResponse{OrderID: createdOrder.ID.String()}

	mockService.On("CreateOrder", ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated).Return(response, nil)
	mockService.On("GetOrder", ctx, createdOrder.ID).Return(createdOrder, nil)

	result, err := uc.CreateOrder(ctx, customerID, courierID)

	assert.NoError(t, err)
	assert.Equal(t, createdOrder, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_CreateOrder_ServiceError(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()

	expectedErr := errors.New("service error")
	mockService.On("CreateOrder", ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated).Return(service.CreateOrderResponse{}, expectedErr)

	result, err := uc.CreateOrder(ctx, customerID, courierID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, models.Order{}, result)
	mockService.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
}

func TestOrderUseCase_CreateOrder_InvalidOrderID(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()

	response := service.CreateOrderResponse{OrderID: "invalid-uuid"}
	mockService.On("CreateOrder", ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated).Return(response, nil)

	result, err := uc.CreateOrder(ctx, customerID, courierID)

	assert.Error(t, err)
	assert.Equal(t, models.Order{}, result)
	mockService.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
}

func TestOrderUseCase_CreateOrderWithItems_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()
	items := []repositoryModels.OrderItemInput{createTestOrderItemInput()}

	createdOrder := createTestOrder()
	response := service.CreateOrderResponse{OrderID: createdOrder.ID.String()}

	mockService.On("CreateOrderWithItems", ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated, items).Return(response, nil)
	mockService.On("GetOrder", ctx, createdOrder.ID).Return(createdOrder, nil)

	result, err := uc.CreateOrderWithItems(ctx, customerID, courierID, items)

	assert.NoError(t, err)
	assert.Equal(t, createdOrder, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_CreateOrderWithItems_InvalidOrderID(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()
	items := []repositoryModels.OrderItemInput{createTestOrderItemInput()}

	response := service.CreateOrderResponse{OrderID: "invalid-uuid"}
	mockService.On("CreateOrderWithItems", ctx, customerID.String(), courierID.String(), models.OrderStatusCustomerCreated, items).Return(response, nil)

	result, err := uc.CreateOrderWithItems(ctx, customerID, courierID, items)

	assert.Error(t, err)
	assert.Equal(t, models.Order{}, result)
	mockService.AssertExpectations(t)
	mockService.AssertNotCalled(t, "GetOrder", mock.Anything, mock.Anything)
}

func TestOrderUseCase_GetOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	expectedOrder := createTestOrder()

	mockService.On("GetOrder", ctx, orderID).Return(expectedOrder, nil)

	result, err := uc.GetOrder(ctx, orderID)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_AcceptOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()
	courierID := uuid.New()
	items := []repositoryModels.OrderItemInput{createTestOrderItemInput()}
	status := models.OrderStatusKitchenAccepted

	expected := repositoryModels.AcceptResult{OrderID: orderID, Status: string(status)}
	mockService.On("AcceptOrder", ctx, orderID.String(), customerID.String(), courierID.String(), items, status).Return(expected, nil)

	result, err := uc.AcceptOrder(ctx, orderID, customerID, courierID, items, status)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_GetOrderStatus_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	expectedStatus := models.OrderStatusDeliveryPending

	mockService.On("GetOrderStatus", ctx, orderID).Return(expectedStatus, nil)

	result, err := uc.GetOrderStatus(ctx, orderID)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_UpdateOrderStatus_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	status := models.OrderStatusOrderCompleted

	mockService.On("UpdateOrderStatus", ctx, orderID, status).Return(nil)

	err := uc.UpdateOrderStatus(ctx, orderID, status)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_CalculateOrderTotal_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	expectedTotal := 42.75

	mockService.On("CalculateOrderTotal", ctx, orderID).Return(expectedTotal, nil)

	result, err := uc.CalculateOrderTotal(ctx, orderID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_GetCustomerWalletAddress_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	customerID := uuid.New()
	expectedAddress := "wallet_123"

	mockService.On("GetCustomerWalletAddress", ctx, customerID).Return(expectedAddress, nil)

	result, err := uc.GetCustomerWalletAddress(ctx, customerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, result)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_AddItemIntoOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	item := createTestOrderItemInput()

	mockService.On("AddItemIntoOrder", ctx, orderID, item).Return(nil)

	err := uc.AddItemIntoOrder(ctx, orderID, item)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestOrderUseCase_RemoveItemFromOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	uc := NewOrderUseCase(mockService)

	ctx := context.Background()
	orderID := uuid.New()
	restaurantItemID := uuid.New()

	mockService.On("RemoveItemFromOrder", ctx, orderID, restaurantItemID).Return(nil)

	err := uc.RemoveItemFromOrder(ctx, orderID, restaurantItemID)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
