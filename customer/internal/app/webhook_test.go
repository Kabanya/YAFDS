package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
	"github.com/Kabanya/YAFDS/pkg/models"
)

func TestMain(m *testing.M) {
	// Initialize logger for tests
	tempFile := "/tmp/test_webhook.log"
	_ = utils.InitFileLogger(tempFile)
	defer utils.CloseLogger()
	defer os.Remove(tempFile)

	os.Exit(m.Run())
}

// mockWebhookUseCase is a mock implementation of WebhookUseCase for testing
type mockWebhookUseCase struct {
	processFunc func(ctx context.Context, orderID uuid.UUID, status string) error
}

func (m *mockWebhookUseCase) ProcessPaymentWebhook(ctx context.Context, orderID uuid.UUID, status string) error {
	if m.processFunc != nil {
		return m.processFunc(ctx, orderID, status)
	}
	return nil
}

func TestNewWebhookHandler(t *testing.T) {
	t.Run("creates handler with custom signature header", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		customHeader := "Stripe-Signature"
		handler := NewWebhookHandler(mockUC, customHeader)

		assert.NotNil(t, handler)
		assert.Equal(t, mockUC, handler.usecase)
		assert.Equal(t, customHeader, handler.signatureHeader)
	})

	t.Run("creates handler with default signature header when empty string provided", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		handler := NewWebhookHandler(mockUC, "")

		assert.NotNil(t, handler)
		assert.Equal(t, "X-Signature", handler.signatureHeader)
	})
}

func TestSetSignatureHeader(t *testing.T) {
	mockUC := &mockWebhookUseCase{}
	handler := NewWebhookHandler(mockUC, "X-Signature")

	newHeader := "Custom-Webhook-Signature"
	handler.SetSignatureHeader(newHeader)

	assert.Equal(t, newHeader, handler.signatureHeader)
}

func TestWebhookHandler_PaymentWebhook_Success(t *testing.T) {
	t.Run("successfully processes valid payment webhook", func(t *testing.T) {
		orderID := uuid.New()
		payload := PaymentWebhookPayload{
			OrderID:   orderID.String(),
			Status:    "paid",
			Amount:    100.50,
			Currency:  "USD",
			PaymentID: "pay_123456",
			Timestamp: 1234567890,
		}

		var processedOrderID uuid.UUID
		var processedStatus string

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				processedOrderID = id
				processedStatus = status
				return nil
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "valid_signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, orderID, processedOrderID)
		assert.Equal(t, "paid", processedStatus)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, orderID.String(), response["order_id"])
	})

	t.Run("successfully processes webhook with minimal required fields", func(t *testing.T) {
		orderID := uuid.New()
		payload := PaymentWebhookPayload{
			OrderID: orderID.String(),
			Status:  "completed",
		}

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				return nil
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature123")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestWebhookHandler_PaymentWebhook_MethodNotAllowed(t *testing.T) {
	mockUC := &mockWebhookUseCase{}
	handler := NewWebhookHandler(mockUC, "X-Signature")

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method+" method is not allowed", func(t *testing.T) {
			req := httptest.NewRequest(method, "/webhook/payment", nil)
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

			var response models.ErrorResponce
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "method not allowed", response.ErrorMessage)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_MissingSignature(t *testing.T) {
	testCases := []struct {
		name            string
		signatureHeader string
	}{
		{
			name:            "default signature header missing",
			signatureHeader: "X-Signature",
		},
		{
			name:            "custom signature header missing",
			signatureHeader: "Stripe-Signature",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUC := &mockWebhookUseCase{}
			handler := NewWebhookHandler(mockUC, tc.signatureHeader)

			payload := PaymentWebhookPayload{
				OrderID: uuid.New().String(),
				Status:  "paid",
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var response models.ErrorResponce
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "missing signature", response.ErrorMessage)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_InvalidPayload(t *testing.T) {
	testCases := []struct {
		name          string
		body          string
		expectedError string
	}{
		{
			name:          "invalid JSON",
			body:          "{invalid json}",
			expectedError: "invalid payload",
		},
		{
			name:          "empty body",
			body:          "",
			expectedError: "invalid payload",
		},
		{
			name:          "non-JSON body",
			body:          "plain text",
			expectedError: "invalid payload",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUC := &mockWebhookUseCase{}
			handler := NewWebhookHandler(mockUC, "X-Signature")

			req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(tc.body))
			req.Header.Set("X-Signature", "signature")
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response models.ErrorResponce
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedError, response.ErrorMessage)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_MissingRequiredFields(t *testing.T) {
	t.Run("missing order_id", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := map[string]string{
			"status": "paid",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "order_id is required", response.ErrorMessage)
	})

	t.Run("empty order_id", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := PaymentWebhookPayload{
			OrderID: "",
			Status:  "paid",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "order_id is required", response.ErrorMessage)
	})

	t.Run("missing status", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := map[string]string{
			"order_id": uuid.New().String(),
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "status is required", response.ErrorMessage)
	})

	t.Run("empty status", func(t *testing.T) {
		mockUC := &mockWebhookUseCase{}
		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := PaymentWebhookPayload{
			OrderID: uuid.New().String(),
			Status:  "",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "status is required", response.ErrorMessage)
	})
}

func TestWebhookHandler_PaymentWebhook_InvalidOrderID(t *testing.T) {
	testCases := []struct {
		name     string
		orderID  string
		expected string
	}{
		{
			name:     "non-UUID format",
			orderID:  "not-a-uuid",
			expected: "invalid order_id format",
		},
		{
			name:     "partial UUID",
			orderID:  "12345678-1234-1234",
			expected: "invalid order_id format",
		},
		{
			name:     "invalid characters",
			orderID:  "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			expected: "invalid order_id format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUC := &mockWebhookUseCase{}
			handler := NewWebhookHandler(mockUC, "X-Signature")

			payload := PaymentWebhookPayload{
				OrderID: tc.orderID,
				Status:  "paid",
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
			req.Header.Set("X-Signature", "signature")
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response models.ErrorResponce
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, response.ErrorMessage)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_UseCaseErrors(t *testing.T) {
	t.Run("order not found returns 409 Conflict", func(t *testing.T) {
		orderID := uuid.New()

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				return sql.ErrNoRows
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := PaymentWebhookPayload{
			OrderID: orderID.String(),
			Status:  "paid",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, sql.ErrNoRows.Error(), response.ErrorMessage)
	})

	t.Run("internal error returns 500 Internal Server Error", func(t *testing.T) {
		orderID := uuid.New()
		customErr := errors.New("database connection failed")

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				return customErr
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := PaymentWebhookPayload{
			OrderID: orderID.String(),
			Status:  "paid",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, customErr.Error(), response.ErrorMessage)
	})

	t.Run("context cancellation returns 500 Internal Server Error", func(t *testing.T) {
		orderID := uuid.New()
		customErr := context.Canceled

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				return customErr
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		payload := PaymentWebhookPayload{
			OrderID: orderID.String(),
			Status:  "paid",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "signature")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response models.ErrorResponce
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, customErr.Error(), response.ErrorMessage)
	})
}

func TestWebhookHandler_PaymentWebhook_CustomSignatureHeader(t *testing.T) {
	testCases := []struct {
		name            string
		signatureHeader string
		headerValue     string
	}{
		{
			name:            "Stripe signature header",
			signatureHeader: "Stripe-Signature",
			headerValue:     "t=123,v1=abc123",
		},
		{
			name:            "PayPal signature header",
			signatureHeader: "PayPal-Auth-Algo",
			headerValue:     "sha256",
		},
		{
			name:            "Custom webhook header",
			signatureHeader: "X-Webhook-Signature-v2",
			headerValue:     "custom_sig_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orderID := uuid.New()
			payload := PaymentWebhookPayload{
				OrderID: orderID.String(),
				Status:  "paid",
			}

			var processed bool
			mockUC := &mockWebhookUseCase{
				processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
					processed = true
					return nil
				},
			}

			handler := NewWebhookHandler(mockUC, tc.signatureHeader)

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
			req.Header.Set(tc.signatureHeader, tc.headerValue)
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.True(t, processed, "usecase should be called with correct signature header")
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_DifferentPaymentStatuses(t *testing.T) {
	statuses := []string{
		"paid",
		"completed",
		"failed",
		"pending",
		"refunded",
		"partially_refunded",
		"cancelled",
	}

	for _, status := range statuses {
		t.Run("status: "+status, func(t *testing.T) {
			orderID := uuid.New()
			var capturedStatus string

			mockUC := &mockWebhookUseCase{
				processFunc: func(ctx context.Context, id uuid.UUID, s string) error {
					capturedStatus = s
					return nil
				},
			}

			handler := NewWebhookHandler(mockUC, "X-Signature")

			payload := PaymentWebhookPayload{
				OrderID: orderID.String(),
				Status:  status,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
			req.Header.Set("X-Signature", "signature")
			w := httptest.NewRecorder()

			handler.PaymentWebhook(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, status, capturedStatus)
		})
	}
}

func TestWebhookHandler_PaymentWebhook_AdditionalFields(t *testing.T) {
	t.Run("payload with all optional fields", func(t *testing.T) {
		orderID := uuid.New()
		payload := PaymentWebhookPayload{
			OrderID:   orderID.String(),
			Status:    "paid",
			Amount:    999.99,
			Currency:  "EUR",
			PaymentID: "pay_example_12345",
			Timestamp: 1704067200,
		}

		mockUC := &mockWebhookUseCase{
			processFunc: func(ctx context.Context, id uuid.UUID, status string) error {
				return nil
			},
		}

		handler := NewWebhookHandler(mockUC, "X-Signature")

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/webhook/payment", strings.NewReader(string(body)))
		req.Header.Set("X-Signature", "sig_with_extra_fields")
		w := httptest.NewRecorder()

		handler.PaymentWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPaymentWebhookPayload_JsonTags(t *testing.T) {
	t.Run("serializes and deserializes correctly", func(t *testing.T) {
		orderID := uuid.New()
		original := PaymentWebhookPayload{
			OrderID:   orderID.String(),
			Status:    "paid",
			Amount:    100.50,
			Currency:  "USD",
			PaymentID: "pay_123",
			Timestamp: 1234567890,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded PaymentWebhookPayload
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, original.OrderID, decoded.OrderID)
		assert.Equal(t, original.Status, decoded.Status)
		assert.Equal(t, original.Amount, decoded.Amount)
		assert.Equal(t, original.Currency, decoded.Currency)
		assert.Equal(t, original.PaymentID, decoded.PaymentID)
		assert.Equal(t, original.Timestamp, decoded.Timestamp)
	})

	t.Run("omitempty fields work correctly", func(t *testing.T) {
		orderID := uuid.New()
		payload := PaymentWebhookPayload{
			OrderID: orderID.String(),
			Status:  "paid",
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		var raw map[string]interface{}
		err = json.Unmarshal(data, &raw)
		require.NoError(t, err)

		assert.Contains(t, raw, "order_id")
		assert.Contains(t, raw, "status")
		assert.NotContains(t, raw, "amount")
		assert.NotContains(t, raw, "currency")
		assert.NotContains(t, raw, "payment_id")
		assert.NotContains(t, raw, "timestamp")
	})
}
