package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Kabanya/YAFDS/pkg/models"
	pkgRepoModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepo - мок репозитория заказов
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) CreateOrderWithItems(ctx context.Context, order models.Order, items []pkgRepoModels.OrderItemInput) (models.Order, error) {
	args := m.Called(ctx, order, items)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) ListOrders(ctx context.Context, filter pkgRepoModels.Filter) ([]models.Order, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepo) AcceptOrder(ctx context.Context, input pkgRepoModels.AcceptInput) (pkgRepoModels.AcceptResult, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(pkgRepoModels.AcceptResult), args.Error(1)
}

func (m *MockOrderRepo) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.OrderStatus), args.Error(1)
}

func (m *MockOrderRepo) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepo) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderRepo) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	args := m.Called(ctx, customerID)
	return args.String(0), args.Error(1)
}

func (m *MockOrderRepo) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item pkgRepoModels.OrderItemInput) error {
	args := m.Called(ctx, orderID, item)
	return args.Error(0)
}

func (m *MockOrderRepo) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	args := m.Called(ctx, orderID, restaurantItemID)
	return args.Error(0)
}

// Вспомогательные функции для тестов
func createTestOrder() models.Order {
	orderID := uuid.New()
	customerID := uuid.New()
	courierID := uuid.New()

	return models.Order{
		ID:         orderID,
		CustomerID: customerID,
		CourierID:  courierID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     models.OrderStatusCustomerCreated,
	}
}

func createTestOrderItemInput() pkgRepoModels.OrderItemInput {
	return pkgRepoModels.OrderItemInput{
		RestaurantItemID: uuid.New(),
		Price:            15.99,
		Quantity:         2,
	}
}

// Тесты для CreateOrder
func TestOrderSvc_CreateOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated

	expectedOrder := createTestOrder()

	mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("models.Order")).Return(expectedOrder, nil)

	// Act
	result, err := service.CreateOrder(ctx, customerID, courierID, status)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID.String(), result.OrderID)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_CreateOrder_InvalidCustomerID(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	invalidCustomerID := "invalid-uuid"
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated

	// Act
	result, err := service.CreateOrder(ctx, invalidCustomerID, courierID, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, CreateOrderResponse{}, result)
	mockRepo.AssertNotCalled(t, "CreateOrder")
}

func TestOrderSvc_CreateOrder_InvalidCourierID(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New().String()
	invalidCourierID := "invalid-uuid"
	status := models.OrderStatusCustomerCreated

	// Act
	result, err := service.CreateOrder(ctx, customerID, invalidCourierID, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, CreateOrderResponse{}, result)
	mockRepo.AssertNotCalled(t, "CreateOrder")
}

func TestOrderSvc_CreateOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated

	expectedError := errors.New("repository error")
	mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("models.Order")).Return(models.Order{}, expectedError)

	// Act
	result, err := service.CreateOrder(ctx, customerID, courierID, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, CreateOrderResponse{}, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для CreateOrderWithItems
func TestOrderSvc_CreateOrderWithItems_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated
	items := []pkgRepoModels.OrderItemInput{
		createTestOrderItemInput(),
		createTestOrderItemInput(),
	}

	expectedOrder := createTestOrder()

	mockRepo.On("CreateOrderWithItems", ctx, mock.AnythingOfType("models.Order"), items).Return(expectedOrder, nil)

	// Act
	result, err := service.CreateOrderWithItems(ctx, customerID, courierID, status, items)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID.String(), result.OrderID)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_CreateOrderWithItems_InvalidCustomerID(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	invalidCustomerID := "invalid-uuid"
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}

	// Act
	result, err := service.CreateOrderWithItems(ctx, invalidCustomerID, courierID, status, items)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, CreateOrderResponse{}, result)
	mockRepo.AssertNotCalled(t, "CreateOrderWithItems")
}

func TestOrderSvc_CreateOrderWithItems_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	status := models.OrderStatusCustomerCreated
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}

	expectedError := errors.New("repository error")
	mockRepo.On("CreateOrderWithItems", ctx, mock.AnythingOfType("models.Order"), items).Return(models.Order{}, expectedError)

	// Act
	result, err := service.CreateOrderWithItems(ctx, customerID, courierID, status, items)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, CreateOrderResponse{}, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для ListOrders
func TestOrderSvc_ListOrders_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	filter := pkgRepoModels.Filter{
		Status: "CUSTOMER_CREATED",
	}

	expectedOrders := []models.Order{
		createTestOrder(),
		createTestOrder(),
	}

	mockRepo.On("ListOrders", ctx, filter).Return(expectedOrders, nil)

	// Act
	result, err := service.ListOrders(ctx, filter)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_ListOrders_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	filter := pkgRepoModels.Filter{}

	expectedError := errors.New("repository error")
	mockRepo.On("ListOrders", ctx, filter).Return([]models.Order{}, expectedError)

	// Act
	result, err := service.ListOrders(ctx, filter)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для GetOrder
func TestOrderSvc_GetOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedOrder := createTestOrder()

	mockRepo.On("GetOrder", ctx, orderID).Return(expectedOrder, nil)

	// Act
	result, err := service.GetOrder(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_GetOrder_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedError := errors.New("order not found")
	mockRepo.On("GetOrder", ctx, orderID).Return(models.Order{}, expectedError)

	// Act
	result, err := service.GetOrder(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.Order{}, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для AcceptOrder
func TestOrderSvc_AcceptOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}
	status := models.OrderStatusKitchenAccepted

	expectedResult := pkgRepoModels.AcceptResult{
		OrderID: uuid.New(),
		Status:  "ACCEPTED",
	}

	expectedOrderID := uuid.MustParse(orderID)
	expectedCustomerID := uuid.MustParse(customerID)
	expectedCourierID := uuid.MustParse(courierID)

	mockRepo.On("AcceptOrder", ctx, mock.MatchedBy(func(input pkgRepoModels.AcceptInput) bool {
		return input.OrderID == expectedOrderID &&
			input.CustomerID == expectedCustomerID &&
			input.CourierID == expectedCourierID &&
			input.Status == status &&
			assert.ObjectsAreEqual(items, input.Items)
	})).Return(expectedResult, nil)

	// Act
	result, err := service.AcceptOrder(ctx, orderID, customerID, courierID, items, status)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_AcceptOrder_InvalidOrderID(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	invalidOrderID := "invalid-uuid"
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}
	status := models.OrderStatusKitchenAccepted

	// Act
	result, err := service.AcceptOrder(ctx, invalidOrderID, customerID, courierID, items, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, pkgRepoModels.AcceptResult{}, result)
	mockRepo.AssertNotCalled(t, "AcceptOrder")
}

func TestOrderSvc_AcceptOrder_InvalidCustomerID(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New().String()
	invalidCustomerID := "invalid-uuid"
	courierID := uuid.New().String()
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}
	status := models.OrderStatusKitchenAccepted

	// Act
	result, err := service.AcceptOrder(ctx, orderID, invalidCustomerID, courierID, items, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, pkgRepoModels.AcceptResult{}, result)
	mockRepo.AssertNotCalled(t, "AcceptOrder")
}

func TestOrderSvc_AcceptOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New().String()
	customerID := uuid.New().String()
	courierID := uuid.New().String()
	items := []pkgRepoModels.OrderItemInput{createTestOrderItemInput()}
	status := models.OrderStatusKitchenAccepted

	expectedError := errors.New("repository error")

	expectedOrderID := uuid.MustParse(orderID)
	expectedCustomerID := uuid.MustParse(customerID)
	expectedCourierID := uuid.MustParse(courierID)

	mockRepo.On("AcceptOrder", ctx, mock.MatchedBy(func(input pkgRepoModels.AcceptInput) bool {
		return input.OrderID == expectedOrderID &&
			input.CustomerID == expectedCustomerID &&
			input.CourierID == expectedCourierID &&
			input.Status == status &&
			assert.ObjectsAreEqual(items, input.Items)
	})).Return(pkgRepoModels.AcceptResult{}, expectedError)

	// Act
	result, err := service.AcceptOrder(ctx, orderID, customerID, courierID, items, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, pkgRepoModels.AcceptResult{}, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для GetOrderStatus
func TestOrderSvc_GetOrderStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedStatus := models.OrderStatusCustomerPaid

	mockRepo.On("GetOrderStatus", ctx, orderID).Return(expectedStatus, nil)

	// Act
	result, err := service.GetOrderStatus(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_GetOrderStatus_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedError := errors.New("repository error")
	mockRepo.On("GetOrderStatus", ctx, orderID).Return(models.OrderStatus(""), expectedError)

	// Act
	result, err := service.GetOrderStatus(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.OrderStatus(""), result)
	mockRepo.AssertExpectations(t)
}

// Тесты для UpdateOrderStatus
func TestOrderSvc_UpdateOrderStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	status := models.OrderStatusOrderCompleted

	mockRepo.On("UpdateOrderStatus", ctx, orderID, status).Return(nil)

	// Act
	err := service.UpdateOrderStatus(ctx, orderID, status)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_UpdateOrderStatus_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	status := models.OrderStatusOrderCompleted

	expectedError := errors.New("repository error")
	mockRepo.On("UpdateOrderStatus", ctx, orderID, status).Return(expectedError)

	// Act
	err := service.UpdateOrderStatus(ctx, orderID, status)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
func TestOrderSvc_CalculateOrderTotal_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedTotal := 42.99

	mockRepo.On("CalculateOrderTotal", ctx, orderID).Return(expectedTotal, nil)

	// Act
	result, err := service.CalculateOrderTotal(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_CalculateOrderTotal_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()

	expectedError := errors.New("repository error")
	mockRepo.On("CalculateOrderTotal", ctx, orderID).Return(0.0, expectedError)

	// Act
	result, err := service.CalculateOrderTotal(ctx, orderID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0.0, result)
	mockRepo.AssertExpectations(t)
}

// Тесты для GetCustomerWalletAddress
func TestOrderSvc_GetCustomerWalletAddress_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New()

	expectedAddress := "0x1234567890abcdef"

	mockRepo.On("GetCustomerWalletAddress", ctx, customerID).Return(expectedAddress, nil)

	// Act
	result, err := service.GetCustomerWalletAddress(ctx, customerID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_GetCustomerWalletAddress_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	customerID := uuid.New()

	expectedError := errors.New("repository error")
	mockRepo.On("GetCustomerWalletAddress", ctx, customerID).Return("", expectedError)

	// Act
	result, err := service.GetCustomerWalletAddress(ctx, customerID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, "", result)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_AddItemIntoOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	item := createTestOrderItemInput()

	mockRepo.On("AddItemIntoOrder", ctx, orderID, item).Return(nil)

	// Act
	err := service.AddItemIntoOrder(ctx, orderID, item)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderSvc_AddItemIntoOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepo)
	service := NewOrderService(mockRepo)

	ctx := context.Background()
	orderID := uuid.New()
	item := createTestOrderItemInput()

	expectedError := errors.New("repository error")
	mockRepo.On("AddItemIntoOrder", ctx, orderID, item).Return(expectedError)

	// Act
	err := service.AddItemIntoOrder(ctx, orderID, item)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
