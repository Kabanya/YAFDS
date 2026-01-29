package service

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/google/uuid"
)

func init() {
	// Initialize logger for tests
	utils.InitFileLogger(os.DevNull)
}

// mockOrderRepo mocks OrderRepo
type mockOrderRepo struct {
	getStatusFunc    func(ctx context.Context, orderID uuid.UUID) (string, error)
	updateStatusFunc func(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus, currentStatus string) error
	updatedStatuses  []orderStatusUpdate
}

type orderStatusUpdate struct {
	OrderID       uuid.UUID
	NewStatus     models.OrderStatus
	CurrentStatus string
}

func (m *mockOrderRepo) GetOrderStatus(ctx context.Context, orderID uuid.UUID) (string, error) {
	if m.getStatusFunc != nil {
		return m.getStatusFunc(ctx, orderID)
	}
	return string(models.OrderStatusCustomerCreated), nil
}

func (m *mockOrderRepo) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus, currentStatus string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, orderID, newStatus, currentStatus)
	}
	m.updatedStatuses = append(m.updatedStatuses, orderStatusUpdate{
		OrderID:       orderID,
		NewStatus:     newStatus,
		CurrentStatus: currentStatus,
	})
	return nil
}

// mockRabbitMQPublisher mocks RabbitMQPublisher
type mockRabbitMQPublisher struct {
	publishFunc func(orderID uuid.UUID) error
	closeFunc   func() error
	published   []uuid.UUID
}

func (m *mockRabbitMQPublisher) PublishPaymentNotification(orderID uuid.UUID) error {
	if m.publishFunc != nil {
		return m.publishFunc(orderID)
	}
	m.published = append(m.published, orderID)
	return nil
}

func (m *mockRabbitMQPublisher) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

// TestValidateSignature_ValidSignature tests signature validation with correct signature
func TestValidateSignature_ValidSignature(t *testing.T) {
	secret := "test_secret"
	body := []byte(`{"order_id": "123", "status": "paid"}`)

	repo := &mockOrderRepo{}
	publisher := &mockRabbitMQPublisher{}
	svc := NewWebhookService(repo, publisher, secret).(*webhookService)

	// Generate valid signature
	signature := svc.generateSignature(body)

	result := svc.ValidateSignature(body, signature)

	if !result {
		t.Error("expected valid signature to return true")
	}
}

// TestValidateSignature_InvalidSignature tests signature validation with incorrect signature
func TestValidateSignature_InvalidSignature(t *testing.T) {
	secret := "test_secret"
	body := []byte(`{"order_id": "123", "status": "paid"}`)

	repo := &mockOrderRepo{}
	publisher := &mockRabbitMQPublisher{}
	svc := NewWebhookService(repo, publisher, secret).(*webhookService)

	result := svc.ValidateSignature(body, "invalid_signature")

	if result {
		t.Error("expected invalid signature to return false")
	}
}

// TestValidateSignature_EmptySignature tests signature validation with empty signature
func TestValidateSignature_EmptySignature(t *testing.T) {
	secret := "test_secret"
	body := []byte(`{"order_id": "123", "status": "paid"}`)

	repo := &mockOrderRepo{}
	publisher := &mockRabbitMQPublisher{}
	svc := NewWebhookService(repo, publisher, secret).(*webhookService)

	result := svc.ValidateSignature(body, "")

	if result {
		t.Error("expected empty signature to return false")
	}
}

// TestValidateSignature_EmptyBody tests signature validation with empty body
func TestValidateSignature_EmptyBody(t *testing.T) {
	secret := "test_secret"
	body := []byte{}

	repo := &mockOrderRepo{}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, secret).(*webhookService)

	signature := svc.generateSignature(body)
	result := svc.ValidateSignature(body, signature)

	if !result {
		t.Error("expected valid signature for empty body to return true")
	}
}

// TestGenerateSignature_DifferentBodies tests signature generation for different bodies
func TestGenerateSignature_DifferentBodies(t *testing.T) {
	secret := "test_secret"
	body1 := []byte(`{"status": "paid"}`)
	body2 := []byte(`{"status": "failed"}`)

	repo := &mockOrderRepo{}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, secret).(*webhookService)

	sig1 := svc.generateSignature(body1)
	sig2 := svc.generateSignature(body2)

	if sig1 == sig2 {
		t.Error("expected different signatures for different bodies")
	}
}

// TestProcessPaymentWebhook_Paid_Success tests successful payment processing
func TestProcessPaymentWebhook_Paid_Success(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(repo.updatedStatuses) != 1 {
		t.Fatalf("expected 1 status update, got %d", len(repo.updatedStatuses))
	}

	if repo.updatedStatuses[0].NewStatus != models.OrderStatusCustomerPaid {
		t.Errorf("expected status %s, got %s", models.OrderStatusCustomerPaid, repo.updatedStatuses[0].NewStatus)
	}
}

// TestProcessPaymentWebhook_Paid_AlreadyPaid tests idempotency - already paid order
func TestProcessPaymentWebhook_Paid_AlreadyPaid(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerPaid), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error for already paid order, got %v", err)
	}

	// Update should not be called for already paid order
	if len(repo.updatedStatuses) != 0 {
		t.Errorf("expected 0 status updates for already paid order, got %d", len(repo.updatedStatuses))
	}
}

// TestProcessPaymentWebhook_Paid_InvalidTransition tests invalid status transition
func TestProcessPaymentWebhook_Paid_InvalidTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus models.OrderStatus
	}{
		{
			name:          "from kitchen accepted",
			currentStatus: models.OrderStatusKitchenAccepted,
		},
		{
			name:          "from kitchen preparing",
			currentStatus: models.OrderStatusKitchenPreparing,
		},
		{
			name:          "from delivery delivering",
			currentStatus: models.OrderStatusDeliveryDelivering,
		},
		{
			name:          "from order completed",
			currentStatus: models.OrderStatusOrderCompleted,
		},
		{
			name:          "from customer cancelled",
			currentStatus: models.OrderStatusCustomerCancelled,
		},
		{
			name:          "from order failed",
			currentStatus: models.OrderStatusOrderFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID := uuid.New()
			status := string(PaymentStatusPaid)

			repo := &mockOrderRepo{
				getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
					return string(tt.currentStatus), nil
				},
			}
			svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

			err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

			if err == nil {
				t.Error("expected error for invalid status transition, got nil")
			}
		})
	}
}

// TestProcessPaymentWebhook_Paid_ConcurrentUpdate tests concurrent update handling
func TestProcessPaymentWebhook_Paid_ConcurrentUpdate(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
		updateStatusFunc: func(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus, currentStatus string) error {
			return sql.ErrNoRows
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected error for concurrent update, got nil")
	}
	if !errors.Is(err, sql.ErrNoRows) && err.Error() != "concurrent update: sql: no rows in result set" {
		t.Errorf("expected concurrent update error, got %v", err)
	}
}

// TestProcessPaymentWebhook_Paid_GetStatusError tests error when getting current status
func TestProcessPaymentWebhook_Paid_GetStatusError(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)
	expectedErr := errors.New("database connection failed")

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return "", expectedErr
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

// TestProcessPaymentWebhook_Paid_UpdateError tests error when updating status
func TestProcessPaymentWebhook_Paid_UpdateError(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)
	expectedErr := errors.New("update failed")

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
		updateStatusFunc: func(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus, currentStatus string) error {
			return expectedErr
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

// TestProcessPaymentWebhook_Failed tests failed payment status
func TestProcessPaymentWebhook_Failed(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusFailed)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error for failed payment, got %v", err)
	}

	// Status should not be updated for failed payments
	if len(repo.updatedStatuses) != 0 {
		t.Errorf("expected 0 status updates for failed payment, got %d", len(repo.updatedStatuses))
	}
}

// TestProcessPaymentWebhook_Pending tests pending payment status
func TestProcessPaymentWebhook_Pending(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPending)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error for pending payment, got %v", err)
	}

	// Status should not be updated for pending payments
	if len(repo.updatedStatuses) != 0 {
		t.Errorf("expected 0 status updates for pending payment, got %d", len(repo.updatedStatuses))
	}
}

// TestProcessPaymentWebhook_UnknownStatus tests unknown payment status
func TestProcessPaymentWebhook_UnknownStatus(t *testing.T) {
	orderID := uuid.New()
	status := "unknown_status"

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected error for unknown status, got nil")
	}
}

// TestProcessPaymentWebhook_EmptyStatus tests empty payment status
func TestProcessPaymentWebhook_EmptyStatus(t *testing.T) {
	orderID := uuid.New()
	status := ""

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected error for empty status, got nil")
	}
}

// TestProcessPaymentWebhook_NilUUID tests handling of nil UUID
func TestProcessPaymentWebhook_NilUUID(t *testing.T) {
	orderID := uuid.Nil
	status := string(PaymentStatusPaid)

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return "", sql.ErrNoRows
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected error for nil UUID, got nil")
	}
}

// TestProcessPaymentWebhook_CancelledContext tests context cancellation handling
func TestProcessPaymentWebhook_CancelledContext(t *testing.T) {
	orderID := uuid.New()
	status := string(PaymentStatusPaid)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	repo := &mockOrderRepo{
		getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
			return string(models.OrderStatusCustomerCreated), nil
		},
	}
	svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

	err := svc.ProcessPaymentWebhook(ctx, orderID, status)

	// Should handle cancelled context
	_ = err
}

// TestProcessPaymentWebhook_CaseSensitiveStatus tests case sensitivity of status
func TestProcessPaymentWebhook_CaseSensitiveStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectError bool
	}{
		{
			name:        "lowercase paid",
			status:      "paid",
			expectError: false,
		},
		{
			name:        "uppercase PAID",
			status:      "PAID",
			expectError: true,
		},
		{
			name:        "mixed case Paid",
			status:      "Paid",
			expectError: true,
		},
		{
			name:        "lowercase failed",
			status:      "failed",
			expectError: false,
		},
		{
			name:        "lowercase pending",
			status:      "pending",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID := uuid.New()

			repo := &mockOrderRepo{
				getStatusFunc: func(ctx context.Context, orderID uuid.UUID) (string, error) {
					return string(models.OrderStatusCustomerCreated), nil
				},
			}
			svc := NewWebhookService(repo, &mockRabbitMQPublisher{}, "secret")

			err := svc.ProcessPaymentWebhook(context.Background(), orderID, tt.status)

			if tt.expectError && err == nil {
				t.Errorf("expected error for status '%s', got nil", tt.status)
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error for status '%s', got %v", tt.status, err)
			}
		})
	}
}

// TestPaymentStatus_Constants tests payment status constants
func TestPaymentStatus_Constants(t *testing.T) {
	tests := []struct {
		status   PaymentStatus
		expected string
	}{
		{PaymentStatusPaid, "paid"},
		{PaymentStatusFailed, "failed"},
		{PaymentStatusPending, "pending"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.status)
		}
	}
}

// TestNewWebhookService tests webhook service creation
func TestNewWebhookService(t *testing.T) {
	repo := &mockOrderRepo{}
	publisher := &mockRabbitMQPublisher{}
	secret := "test_secret"

	svc := NewWebhookService(repo, publisher, secret)

	if svc == nil {
		t.Error("expected non-nil service")
	}

	typedSvc, ok := svc.(*webhookService)
	if !ok {
		t.Error("expected *webhookService type")
	}

	if typedSvc.secret != secret {
		t.Errorf("expected secret %s, got %s", secret, typedSvc.secret)
	}

	if typedSvc.orderRepo != repo {
		t.Error("expected order repo to be set")
	}
}
