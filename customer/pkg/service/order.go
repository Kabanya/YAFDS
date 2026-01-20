package service

// НЕ ЗНАЕТ ПРО HTTP .
// Создается, есть репу, логеры.
// есть func create order
import (
	order "customer/pkg/repository"
)

type Repository = order.Repository
type Filter = order.Filter

type createOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status"`
}

type acceptOrderRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type createOrderResponce struct {
	OrderID string `json:"order_id"`
}

type createOrderErrorResponse struct {
	Error string `json:"error"`
}
