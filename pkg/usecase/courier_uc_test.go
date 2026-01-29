package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Kabanya/YAFDS/pkg/models"
	. "github.com/Kabanya/YAFDS/pkg/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCourierService struct {
	mock.Mock
}

func (m *MockCourierService) ListCouriers(ctx context.Context) ([]models.Courier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Courier), args.Error(1)
}

func createTestCourier() models.Courier {
	return models.Courier{
		ID:   uuid.New(),
		Name: "Test Courier",
	}
}

func TestCourierUseCase_ListCouriers_Success(t *testing.T) {
	mockService := new(MockCourierService)
	uc := NewCourierUseCase(mockService)

	ctx := context.Background()
	courier1 := createTestCourier()
	courier2 := createTestCourier()
	courier2.Name = "Test Courier 2"

	expectedCouriers := []models.Courier{courier1, courier2}

	mockService.On("ListCouriers", ctx).Return(expectedCouriers, nil)

	result, err := uc.ListCouriers(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedCouriers, result)
	assert.Len(t, result, 2)
	mockService.AssertExpectations(t)
}

func TestCourierUseCase_ListCouriers_EmptyList(t *testing.T) {
	mockService := new(MockCourierService)
	uc := NewCourierUseCase(mockService)

	ctx := context.Background()
	expectedCouriers := []models.Courier{}

	mockService.On("ListCouriers", ctx).Return(expectedCouriers, nil)

	result, err := uc.ListCouriers(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedCouriers, result)
	assert.Len(t, result, 0)
	mockService.AssertExpectations(t)
}

func TestCourierUseCase_ListCouriers_ServiceError(t *testing.T) {
	mockService := new(MockCourierService)
	uc := NewCourierUseCase(mockService)

	ctx := context.Background()
	expectedError := errors.New("database connection error")

	mockService.On("ListCouriers", ctx).Return([]models.Courier{}, expectedError)

	result, err := uc.ListCouriers(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}
