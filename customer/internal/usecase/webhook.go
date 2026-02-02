package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/google/uuid"
)

// WebhookUseCase defines the interface for webhook operations
type WebhookUseCase interface {
	ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error
}

type webhookUseCase struct {
	service WebhookService
}

type WebhookService interface {
	ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error
}

func NewWebhookUseCase(service WebhookService) WebhookUseCase {
	return &webhookUseCase{
		service: service,
	}
}

func (u *webhookUseCase) ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error {
	return u.service.ProcessPaymentWebhook(ctx, orderID, status)
}

func WithTimeout(inner WebhookUseCase, timeout time.Duration) WebhookUseCase {
	return &timeoutWebhookUseCase{inner, timeout}
}

type timeoutWebhookUseCase struct {
	inner   WebhookUseCase
	timeout time.Duration
}

func (t *timeoutWebhookUseCase) ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- t.inner.ProcessPaymentWebhook(ctx, orderID, status)
	}()

	select {
	case err := <-errChan:
		if errors.Is(err, context.DeadlineExceeded) {
			logger, _ := utils.Logger()
			logger.Printf("Webhook processing timeout for order %s", orderID)
			return context.DeadlineExceeded
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
