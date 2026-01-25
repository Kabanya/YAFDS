// создать enum для статусов заказа
package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusCustomerCreated    OrderStatus = "CUSTOMER_CREATED"
	OrderStatusCustomerPaid       OrderStatus = "CUSTOMER_PAID"
	OrderStatusCustomerCancelled  OrderStatus = "CUSTOMER_CANCELLED"
	OrderStatusKitchenAccepted    OrderStatus = "KITCHEN_ACCEPTED"
	OrderStatusKitchenDenied      OrderStatus = "KITCHEN_DENIED"
	OrderStatusKitchenPreparing   OrderStatus = "KITCHEN_PREPARING"
	OrderStatusCourierRefunded    OrderStatus = "COURIER_REFUNDED"
	OrderStatusDeliveryPending    OrderStatus = "DELIVERY_PENDING"
	OrderStatusDeliveryPicking    OrderStatus = "DELIVERY_PICKING"
	OrderStatusDeliveryDenied     OrderStatus = "DELIVERY_DENIED"
	OrderStatusDeliveryRefunded   OrderStatus = "DELIVERY_REFUNDED"
	OrderStatusDeliveryDelivering OrderStatus = "DELIVERY_DELIVERING"
	OrderStatusOrderCompleted     OrderStatus = "ORDER_COMPLETED"
)

type ErrorResponce struct {
	ErrorMessage string `json:"error_message"`
}

type Order struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	CourierID  uuid.UUID `json:"courier_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
}

type MenuItem struct {
	OrderItemID  uuid.UUID `json:"order_item_id" db:"order_item_id"`
	RestaurantID uuid.UUID `json:"restaurant_id" db:"restaurant_id"`
	Name         string    `json:"name" db:"name"`
	Price        float64   `json:"price" db:"price"`
	Quantity     int       `json:"quantity" db:"quantity"`
	Description  string    `json:"description" db:"description"`
}
