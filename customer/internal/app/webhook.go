package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/google/uuid"
)

// PaymentWebhookPayload represents the incoming payment webhook data
type PaymentWebhookPayload struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	// Additional fields from payment provider
	Amount    float64 `json:"amount,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	PaymentID string  `json:"payment_id,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// WebhookHandler handles payment webhook HTTP endpoints
type WebhookHandler struct {
	usecase         WebhookUseCase
	signatureHeader string // e.g., "X-Signature", "Stripe-Signature"
}

type WebhookUseCase interface {
	ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error
}

// NewWebhookHandler creates a new payment webhook handler
func NewWebhookHandler(usecase WebhookUseCase, signatureHeader string) *WebhookHandler {
	if signatureHeader == "" {
		signatureHeader = "X-Signature"
	}
	return &WebhookHandler{
		usecase:         usecase,
		signatureHeader: signatureHeader,
	}
}

// PaymentWebhook handles incoming payment webhooks
func (wh *WebhookHandler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	logger, _ := utils.Logger()
	logger.Println("Payment webhook called")

	// Only accept POST requests
	if r.Method != http.MethodPost {
		utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the raw body for signature validation
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Failed to read request body: %v", err)
		utils.WriteError(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate signature
	signature := r.Header.Get(wh.signatureHeader)
	if signature == "" {
		logger.Printf("Missing signature header: %s", wh.signatureHeader)
		utils.WriteError(w, "missing signature", http.StatusUnauthorized)
		return
	}

	// Parse webhook payload
	var payload PaymentWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		logger.Printf("Failed to parse webhook payload: %v", err)
		utils.WriteError(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if payload.OrderID == "" {
		utils.WriteError(w, "order_id is required", http.StatusBadRequest)
		return
	}

	if payload.Status == "" {
		utils.WriteError(w, "status is required", http.StatusBadRequest)
		return
	}

	// Parse order ID
	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		logger.Printf("Invalid order ID format: %v", err)
		utils.WriteError(w, "invalid order_id format", http.StatusBadRequest)
		return
	}

	// Process payment webhook via usecase
	if err := wh.usecase.ProcessPaymentWebhook(r.Context(), orderID, payload.Status); err != nil {
		logger.Printf("Failed to process payment webhook: %v", err)

		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusConflict
		}

		utils.WriteError(w, err.Error(), statusCode)
		return
	}

	logger.Printf("Order %s payment webhook processed successfully", orderID)
	utils.WriteJSON(w, map[string]string{
		"status":   "success",
		"order_id": payload.OrderID,
	}, http.StatusOK)
}

// SetSignatureHeader allows customizing the signature header name
func (wh *WebhookHandler) SetSignatureHeader(header string) {
	wh.signatureHeader = header
}
