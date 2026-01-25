package service

// НЕ ЗНАЕТ ПРО HTTP .
// Создается, есть репу, логеры.
// есть func create order
import (
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
)

type Repository = repositoryModels.Order
type Filter = repositoryModels.Filter

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
