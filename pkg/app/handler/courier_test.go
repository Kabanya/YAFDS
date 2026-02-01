package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/Kabanya/YAFDS/pkg/app/handler"
	"github.com/Kabanya/YAFDS/pkg/models"
	"github.com/google/uuid"
)

type mockCourierUC struct {
	res []models.Courier
	err error
}

func (m *mockCourierUC) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	return m.res, m.err
}

func TestNewCouriersHandler_Success(t *testing.T) {
	id := uuid.New()
	expected := []models.Courier{
		{ID: id, Name: "Ivan", TransportType: "bike", IsActive: true},
	}
	mock := &mockCourierUC{res: expected}
	h := NewCouriersHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var got []models.Courier
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("body mismatch: got %#v, want %#v", got, expected)
	}
}

func TestNewCouriersHandler_Error(t *testing.T) {
	mock := &mockCourierUC{res: nil, err: errors.New("db fail")}
	h := NewCouriersHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rr.Code)
	}

	var errResp struct {
		ErrorMessage string `json:"error_message"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.ErrorMessage != "db fail" {
		t.Fatalf("unexpected error message: %q", errResp.ErrorMessage)
	}
}
