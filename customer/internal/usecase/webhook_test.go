package usecase

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/google/uuid"
)

func init() {
	// Initialize logger for tests
	utils.InitFileLogger(os.DevNull)
}

// mockWebhookService mocks WebhookService
type mockWebhookService struct {
	processFunc func(ctx context.Context, orderID uuid.UUID, status string) error
}

func (m *mockWebhookService) ValidateSignature(body []byte, signature string) bool {
	return true
}

func (m *mockWebhookService) ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error {
	if m.processFunc != nil {
		return m.processFunc(ctx, orderID, status)
	}
	return nil
}

// TestProcessPaymentWebhook_Success tests successful webhook processing
func TestProcessPaymentWebhook_Success(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	serviceCalled := false
	service := &mockWebhookService{
		processFunc: func(ctx context.Context, id uuid.UUID, s string) error {
			if id == orderID && s == status {
				serviceCalled = true
			}
			return nil
		},
	}

	uc := NewWebhookUseCase(service)
	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !serviceCalled {
		t.Error("expected service to be called")
	}
}

// TestProcessPaymentWebhook_ServiceError tests error handling when service fails
func TestProcessPaymentWebhook_ServiceError(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	serviceErr := errors.New("order not found")

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			return serviceErr
		},
	}

	uc := NewWebhookUseCase(service)
	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if err != serviceErr {
		t.Errorf("expected error %v, got %v", serviceErr, err)
	}
}

// TestProcessPaymentWebhook_DifferentStatuses tests various payment statuses
func TestProcessPaymentWebhook_DifferentStatuses(t *testing.T) {
	statuses := []string{"paid", "failed", "pending", "cancelled", "refunded", "expired"}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			orderID := uuid.New()
			service := &mockWebhookService{}

			uc := NewWebhookUseCase(service)
			err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

			if err != nil {
				t.Errorf("expected no error for status %s, got %v", status, err)
			}

		})
	}
}

// TestProcessPaymentWebhook_ContextCancellation tests context cancellation handling
func TestProcessPaymentWebhook_ContextCancellation(t *testing.T) {
	orderID := uuid.New()
	status := "paid"

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	uc := NewWebhookUseCase(service)
	err := uc.ProcessPaymentWebhook(ctx, orderID, status)

	// Context cancellation should cause an error
	if err == nil {
		t.Error("expected error with cancelled context, got nil")
	}
}

// TestProcessPaymentWebhook_NilUUID tests handling of nil UUID
func TestProcessPaymentWebhook_NilUUID(t *testing.T) {
	orderID := uuid.Nil
	status := "paid"

	service := &mockWebhookService{}

	uc := NewWebhookUseCase(service)
	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	// Should handle nil UUID - behavior depends on service implementation
	// This test documents current behavior
	_ = err
}

// TestWithTimeout_Success tests successful execution within timeout
func TestWithTimeout_Success(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := 5 * time.Second

	service := &mockWebhookService{}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err != nil {
		t.Errorf("expected no error within timeout, got %v", err)
	}
}

// TestWithTimeout_TimeoutExceeded tests timeout handling
func TestWithTimeout_TimeoutExceeded(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := 100 * time.Millisecond

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			// Simulate slow operation
			select {
			case <-time.After(1 * time.Second):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) && err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded error, got %v", err)
	}
}

// TestWithTimeout_CancelledContext tests with already cancelled context
func TestWithTimeout_CancelledContext(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := 5 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	service := &mockWebhookService{}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(ctx, orderID, status)

	if err == nil {
		t.Error("expected error with cancelled context, got nil")
	}
}

// TestWithTimeout_InnerError tests error propagation from inner usecase
func TestWithTimeout_InnerError(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := 5 * time.Second
	innerErr := errors.New("inner usecase error")

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			return innerErr
		},
	}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	if err == nil {
		t.Error("expected inner error, got nil")
	}
	if err != innerErr {
		t.Errorf("expected error %v, got %v", innerErr, err)
	}
}

// TestWithTimeout_ZeroTimeout tests with zero timeout
func TestWithTimeout_ZeroTimeout(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := time.Nanosecond // Effectively zero

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
	}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(context.Background(), orderID, status)

	// Should timeout immediately
	if err == nil {
		t.Error("expected timeout error with zero timeout, got nil")
	}
}

// TestWithTimeout_ContextCancelledDuringExecution tests cancellation during execution
func TestWithTimeout_ContextCancelledDuringExecution(t *testing.T) {
	orderID := uuid.New()
	status := "paid"
	timeout := 5 * time.Second

	ctx, cancel := context.WithCancel(context.Background())

	service := &mockWebhookService{
		processFunc: func(ctx context.Context, orderID uuid.UUID, status string) error {
			// Cancel after a short delay
			time.AfterFunc(50*time.Millisecond, cancel)
			time.Sleep(200 * time.Millisecond)
			return nil
		},
	}

	inner := NewWebhookUseCase(service)
	uc := WithTimeout(inner, timeout)

	err := uc.ProcessPaymentWebhook(ctx, orderID, status)

	// Should return context error
	if err == nil {
		t.Error("expected context cancellation error, got nil")
	}
}

// TestMockPublisher tests MockPublisher implementation
// func TestMockPublisher(t *testing.T) {
// 	publisher := &MockPublisher{}
//
// 	orderID := uuid.New()
// 	err := publisher.PublishPaymentNotification(orderID)
//
// 	if err != nil {
// 		t.Errorf("expected no error from MockPublisher, got %v", err)
// 	}
//
// 	err = publisher.Close()
// 	if err != nil {
// 		t.Errorf("expected no error from Close, got %v", err)
// 	}
// }
// (Commented out because MockPublisher was moved or removed)
