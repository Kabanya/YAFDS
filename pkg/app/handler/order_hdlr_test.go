package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kabanya/YAFDS/pkg/app/client"
	"github.com/Kabanya/YAFDS/pkg/app/handler"
	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) CreateOrder(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID) (models.Order, error) {
	args := m.Called(ctx, customerID, courierID)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderUseCase) CreateOrderWithItems(ctx context.Context, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput) (models.Order, error) {
	args := m.Called(ctx, customerID, courierID, items)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderUseCase) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderUseCase) AcceptOrder(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID, courierID uuid.UUID, items []repositoryModels.OrderItemInput, status models.OrderStatus) (repositoryModels.AcceptResult, error) {
	args := m.Called(ctx, orderID, customerID, courierID, items, status)
	return args.Get(0).(repositoryModels.AcceptResult), args.Error(1)
}

func (m *MockOrderUseCase) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (models.OrderStatus, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.OrderStatus), args.Error(1)
}

func (m *MockOrderUseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderUseCase) CalculateOrderTotal(ctx context.Context, orderID uuid.UUID) (float64, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderUseCase) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	args := m.Called(ctx, customerID)
	return args.String(0), args.Error(1)
}

func (m *MockOrderUseCase) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
	args := m.Called(ctx, orderID, item)
	return args.Error(0)
}

func (m *MockOrderUseCase) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	args := m.Called(ctx, orderID, restaurantItemID)
	return args.Error(0)
}

func (m *MockOrderUseCase) PayOrder(ctx context.Context, orderID uuid.UUID, walletClient client.WalletClient) error {
	args := m.Called(ctx, orderID, walletClient)
	return args.Error(0)
}

type MockOrderRepo struct {
	mock.Mock
}

// PayOrder implements [models.OrderRepo].
func (m *MockOrderRepo) PayOrder(ctx context.Context, orderID uuid.UUID) error {
	panic("unimplemented")
}

func (m *MockOrderRepo) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) CreateOrderWithItems(ctx context.Context, order models.Order, items []repositoryModels.OrderItemInput) (models.Order, error) {
	args := m.Called(ctx, order, items)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *MockOrderRepo) ListOrders(ctx context.Context, filter repositoryModels.Filter) ([]models.Order, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *MockOrderRepo) AcceptOrder(ctx context.Context, input repositoryModels.AcceptInput) (repositoryModels.AcceptResult, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(repositoryModels.AcceptResult), args.Error(1)
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

func (m *MockOrderRepo) AddItemIntoOrder(ctx context.Context, orderID uuid.UUID, item repositoryModels.OrderItemInput) error {
	args := m.Called(ctx, orderID, item)
	return args.Error(0)
}

func (m *MockOrderRepo) RemoveItemFromOrder(ctx context.Context, orderID uuid.UUID, restaurantItemID uuid.UUID) error {
	args := m.Called(ctx, orderID, restaurantItemID)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	customerID := uuid.New()
	courierID := uuid.New()
	orderID := uuid.New()

	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
		})
		r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		mockUC.On("CreateOrder", mock.Anything, customerID, courierID).Return(models.Order{ID: orderID}, nil).Once()

		h.CreateOrder(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp models.Order
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, orderID, resp.ID)
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		h.CreateOrder(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("invalid customer_id", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"customer_id": "invalid",
			"courier_id":  courierID.String(),
		})
		r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		h.CreateOrder(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
		})
		r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		mockUC.On("CreateOrder", mock.Anything, customerID, courierID).Return(models.Order{}, errors.New("error")).Once()

		h.CreateOrder(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetOrder(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("GetOrder", mock.Anything, orderID).Return(models.Order{ID: orderID}, nil).Once()

		h.GetOrder(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp models.Order
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, orderID, resp.ID)
	})

	t.Run("missing order_id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders/", nil)
		w := httptest.NewRecorder()

		h.GetOrder(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewListHandler(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	h := handler.NewListHandler(mockRepo)

	customerID := uuid.New()

	t.Run("success with customer_id filter", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders?customer_id="+customerID.String(), nil)
		w := httptest.NewRecorder()

		mockRepo.On("ListOrders", mock.Anything, mock.MatchedBy(func(f repositoryModels.Filter) bool {
			return f.CustomerID != nil && *f.CustomerID == customerID
		})).Return([]models.Order{{ID: uuid.New()}}, nil).Once()

		h(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success with courier_id and status filters", func(t *testing.T) {
		courierID := uuid.New()
		status := "COMPLETED"
		r := httptest.NewRequest(http.MethodGet, "/orders?courier_id="+courierID.String()+"&status="+status, nil)
		w := httptest.NewRecorder()

		mockRepo.On("ListOrders", mock.Anything, mock.MatchedBy(func(f repositoryModels.Filter) bool {
			return f.CourierID != nil && *f.CourierID == courierID && f.Status == status
		})).Return([]models.Order{}, nil).Once()

		h(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "/orders", nil)
		w := httptest.NewRecorder()

		h(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("empty list", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		mockRepo.On("ListOrders", mock.Anything, mock.Anything).Return([]models.Order(nil), nil).Once()

		h(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "[]\n", w.Body.String())
	})
}

func TestUpdateOrderStatus(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	status := models.OrderStatusKitchenPreparing

	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{
			"status": string(status),
		})
		r := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/status", bytes.NewBuffer(reqBody))
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("UpdateOrderStatus", mock.Anything, orderID, status).Return(nil).Once()

		h.UpdateOrderStatus(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAddItemIntoOrder(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	item := repositoryModels.OrderItemInput{
		RestaurantItemID: uuid.New(),
		Price:            10.5,
		Quantity:         2,
	}

	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"restaurant_item_id": item.RestaurantItemID.String(),
			"price":              item.Price,
			"quantity":           item.Quantity,
		})
		r := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/items", bytes.NewBuffer(reqBody))
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("AddItemIntoOrder", mock.Anything, orderID, item).Return(nil).Once()

		h.AddItemIntoOrder(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid restaurant_item_id", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"restaurant_item_id": "invalid",
			"price":              10.0,
			"quantity":           1,
		})
		r := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/items", bytes.NewBuffer(reqBody))
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		h.AddItemIntoOrder(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateOrderWithItems(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	customerID := uuid.New()
	courierID := uuid.New()
	orderID := uuid.New()
	items := []repositoryModels.OrderItemInput{
		{RestaurantItemID: uuid.New(), Price: 10.0, Quantity: 1},
	}

	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
			"items":       items,
		})
		r := httptest.NewRequest(http.MethodPost, "/orders/with-items", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		mockUC.On("CreateOrderWithItems", mock.Anything, customerID, courierID, items).Return(models.Order{ID: orderID}, nil).Once()

		h.CreateOrderWithItems(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestAcceptOrder(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	customerID := uuid.New()
	courierID := uuid.New()
	items := []repositoryModels.OrderItemInput{
		{RestaurantItemID: uuid.New(), Price: 5.0, Quantity: 2},
	}
	status := models.OrderStatusKitchenAccepted

	t.Run("success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
			"items":       items,
			"status":      status,
		})
		r := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/accept", bytes.NewBuffer(reqBody))
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("AcceptOrder", mock.Anything, orderID, customerID, courierID, items, status).Return(repositoryModels.AcceptResult{OrderID: orderID, Status: string(status)}, nil).Once()

		h.AcceptOrder(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetOrderStatus(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	status := models.OrderStatusOrderCompleted

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String()+"/status", nil)
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("GetOrderStatus", mock.Anything, orderID).Return(status, nil).Once()

		h.GetOrderStatus(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, string(status), resp["status"])
	})
}

func TestCalculateOrderTotal(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	total := 42.5

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String()+"/total", nil)
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("CalculateOrderTotal", mock.Anything, orderID).Return(total, nil).Once()

		h.CalculateOrderTotal(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, total, resp["total"])
	})
}

func TestGetCustomerWalletAddress(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	customerID := uuid.New()
	addr := "0x123456789"

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/customers/"+customerID.String()+"/wallet?customer_id="+customerID.String(), nil)
		w := httptest.NewRecorder()

		mockUC.On("GetCustomerWalletAddress", mock.Anything, customerID).Return(addr, nil).Once()

		h.GetCustomerWalletAddress(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, addr, resp["wallet_address"])
	})

	t.Run("missing customer_id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/customers/wallet", nil)
		w := httptest.NewRecorder()

		h.GetCustomerWalletAddress(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRemoveItemFromOrder(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	h := handler.NewOrderHandler(mockUC)

	orderID := uuid.New()
	restaurantItemID := uuid.New()

	t.Run("success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String()+"/items?restaurant_item_id="+restaurantItemID.String(), nil)
		r.SetPathValue("order_id", orderID.String())
		w := httptest.NewRecorder()

		mockUC.On("RemoveItemFromOrder", mock.Anything, orderID, restaurantItemID).Return(nil).Once()

		h.RemoveItemFromOrder(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestOrdersHandler(t *testing.T) {
	mockUC := new(MockOrderUseCase)
	mockRepo := new(MockOrderRepo)
	h := handler.NewOrderHandler(mockUC)

	t.Run("GET calls listHandler", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		mockRepo.On("ListOrders", mock.Anything, mock.Anything).Return([]models.Order{}, nil).Once()

		hdlr := h.OrdersHandler(mockRepo)
		hdlr(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST calls CreateOrder", func(t *testing.T) {
		customerID := uuid.New()
		courierID := uuid.New()
		reqBody, _ := json.Marshal(map[string]string{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
		})
		r := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		mockUC.On("CreateOrder", mock.Anything, customerID, courierID).Return(models.Order{}, nil).Once()

		hdlr := h.OrdersHandler(mockRepo)
		hdlr(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid method", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPut, "/orders", nil)
		w := httptest.NewRecorder()

		hdlr := h.OrdersHandler(mockRepo)
		hdlr(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
