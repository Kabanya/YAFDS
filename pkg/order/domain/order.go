package domain

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
	OrderStatusOrderFailed        OrderStatus = "ORDER_FAILED"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	CourierID  uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     OrderStatus
}

type MenuItem struct {
	OrderItemID      uuid.UUID
	RestaurantItemID uuid.UUID
	Name             string
	Price            float64
	Quantity         int
	Description      string
}
