package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/google/uuid"
)

// PaymentStatus represents payment statuses from providers
type PaymentStatus string

const (
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)

// WebhookService defines the interface for payment webhook business logic
type WebhookService interface {
	ValidateSignature(body []byte, signature string) bool
	ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error
}

type RabbitMQPublisher interface {
	PublishPaymentNotification(orderID uuid.UUID) error
	Close() error
}

type webhookService struct {
	orderRepo OrderRepo
	publisher RabbitMQPublisher
	secret    string
}

type OrderRepo interface {
	GetOrderStatus(ctx context.Context, orderID uuid.UUID) (string, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus, currentStatus string) error
}

// NewWebhookService creates a new webhook service
func NewWebhookService(orderRepo OrderRepo, publisher RabbitMQPublisher, webhookSecret string) WebhookService {
	return &webhookService{
		orderRepo: orderRepo,
		publisher: publisher,
		secret:    webhookSecret,
	}
}

// ValidateSignature validates the HMAC signature from the payment provider
func (s *webhookService) ValidateSignature(body []byte, signature string) bool {
	logger, _ := utils.Logger()

	expectedSignature := s.generateSignature(body)

	// Constant-time comparison to prevent timing attacks
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		logger.Printf("Signature mismatch: expected %s, got %s", expectedSignature, signature)
		return false
	}

	return true
}

// generateSignature generates HMAC-SHA256 signature
func (s *webhookService) generateSignature(body []byte) string {
	h := hmac.New(sha256.New, []byte(s.secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// ProcessPaymentWebhook processes the payment webhook with transaction safety
func (s *webhookService) ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error {
	logger, _ := utils.Logger()

	// Get current order status
	currentStatus, err := s.orderRepo.GetOrderStatus(ctx, orderID)
	if err != nil {
		return err
	}

	logger.Printf("Current order status: %s, webhook status: %s", currentStatus, status)

	// Validate status transition for paid payments
	paymentStatus := PaymentStatus(status)
	if paymentStatus == PaymentStatusPaid {
		if currentStatus == string(models.OrderStatusCustomerPaid) {
			logger.Printf("Order already paid, ignoring duplicate webhook")
			// Idempotency: return success if already in target state
			return nil
		}
		if currentStatus != string(models.OrderStatusCustomerCreated) {
			logger.Printf("Invalid status transition from %s to CUSTOMER_PAID", currentStatus)
			return fmt.Errorf("invalid status transition from %s to CUSTOMER_PAID", currentStatus)
		}
	}

	// Determine new status based on payment status
	var newStatus models.OrderStatus
	switch paymentStatus {
	case PaymentStatusPaid:
		newStatus = models.OrderStatusCustomerPaid
	case PaymentStatusFailed:
		// For failed payments, acknowledge but don't change status
		logger.Printf("Payment failed for order %s", orderID)
		return nil
	case PaymentStatusPending:
		// Payment still processing, no status change
		logger.Printf("Payment pending for order %s", orderID)
		return nil
	default:
		logger.Printf("Unknown payment status: %s", status)
		return fmt.Errorf("unknown payment status: %s", status)
	}

	// Update order status
	if err := s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatus, currentStatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Printf("Order status changed during processing, concurrent update detected")
			return fmt.Errorf("concurrent update: %w", err)
		}
		return err
	}

	// Publish to RabbitMQ for push notification if status changed to PAID
	if newStatus == models.OrderStatusCustomerPaid && s.publisher != nil {
		if err := s.publisher.PublishPaymentNotification(orderID); err != nil {
			logger.Printf("Failed to publish payment notification: %v", err)
			// Note: Order is already updated, so this is non-critical
		} else {
			logger.Printf("Payment notification published for order %s", orderID)
		}
	}

	logger.Printf("Order %s payment webhook processed successfully", orderID)
	return nil
}
