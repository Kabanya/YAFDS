package app

import (
	"context"

	"customer/pkg/repository"

	"restaurant/pkg/models"

	"github.com/google/uuid"
)

// на это уровне могут быть
// клиенты / ручки /

// ПИШЕМ хендлер создания ордера
// всё что он делает -- вызывает usecase
//
// usecase нет напрямую доступа к репозиторию, поэтому он вызывает сервис

// Type aliases from repository
// type Repository = repository.Repository // антипаттерн

// usecase -- сетка зависимостей

type Filter = repository.Filter
type Order = repository.Order

// Error aliases from repository
var (
	ErrCustomerNotFound = repository.ErrCustomerNotFound
	ErrCourierNotFound  = repository.ErrCourierNotFound
)

type App struct {
	repo Repository
}

type errorResponse struct {
	Error string `json:"error"`
}

type courierResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type restaurantResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createRequest struct {
	CustomerID   string                   `json:"customer_id"`
	CourierID    string                   `json:"courier_id"`
	RestaurantID string                   `json:"restaurant_id"`
	Status       string                   `json:"status"`
	Items        []createOrderItemRequest `json:"items"`
}

type createOrderItemRequest struct {
	RestaurantItemID string `json:"restaurant_item_id"`
	Quantity         int    `json:"quantity"`
}

type acceptOrderItemRequest struct {
	RestaurantItemID string  `json:"restaurant_item_id"`
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
}

type acceptOrderRequest struct {
	CustomerID   string                   `json:"customer_id"`
	CourierID    string                   `json:"courier_id"`
	RestaurantID string                   `json:"restaurant_id"`
	Items        []acceptOrderItemRequest `json:"items"`
}

type addOrderItemRequest struct {
	RestaurantID     string `json:"restaurant_id"`
	RestaurantItemID string `json:"restaurant_item_id"`
	Quantity         int    `json:"quantity"`
}

type payOrderRequest struct {
	CustomerID string `json:"customer_id"`
}

type menuItemResponse struct {
	OrderItemID  uuid.UUID `json:"order_item_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
}

type RestaurantMenuClient interface {
	GetMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuItem, error)
}

const itemNotAvailableError = "ITEM_NOT_AVAILABLE"
