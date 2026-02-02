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

type MockRestaurantService struct {
	mock.Mock
}

func (m *MockRestaurantService) ListRestaurants(ctx context.Context) ([]models.Restaurant, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Restaurant), args.Error(1)
}

func (m *MockRestaurantService) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]models.RestaurantMenuItem, error) {
	args := m.Called(ctx, restaurantID)
	return args.Get(0).([]models.RestaurantMenuItem), args.Error(1)
}

func createTestRestaurant() models.Restaurant {
	return models.Restaurant{
		ID:   uuid.New(),
		Name: "Test Restaurant",
	}
}

func createTestMenuItem(restaurantID uuid.UUID) models.RestaurantMenuItem {
	return models.RestaurantMenuItem{
		ID:           uuid.New(),
		RestaurantID: restaurantID,
		Name:         "Test Menu Item",
		Price:        15.99,
	}
}

func TestRestaurantUseCase_ListRestaurants_Success(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	restaurant1 := createTestRestaurant()
	restaurant2 := createTestRestaurant()
	restaurant2.Name = "Test Restaurant 2"

	expectedRestaurants := []models.Restaurant{restaurant1, restaurant2}

	mockService.On("ListRestaurants", ctx).Return(expectedRestaurants, nil)

	result, err := uc.ListRestaurants(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedRestaurants, result)
	assert.Len(t, result, 2)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_ListRestaurants_EmptyList(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	expectedRestaurants := []models.Restaurant{}

	mockService.On("ListRestaurants", ctx).Return(expectedRestaurants, nil)

	result, err := uc.ListRestaurants(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedRestaurants, result)
	assert.Len(t, result, 0)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_ListRestaurants_ServiceError(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	expectedError := errors.New("database connection error")

	mockService.On("ListRestaurants", ctx).Return([]models.Restaurant{}, expectedError)

	result, err := uc.ListRestaurants(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_GetMenu_Success(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	restaurantID := uuid.New()
	item1 := createTestMenuItem(restaurantID)
	item2 := createTestMenuItem(restaurantID)
	item2.Name = "Test Menu Item 2"
	item2.Price = 25.50

	expectedMenu := []models.RestaurantMenuItem{item1, item2}

	mockService.On("GetMenu", ctx, restaurantID).Return(expectedMenu, nil)

	result, err := uc.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenu, result)
	assert.Len(t, result, 2)
	assert.Equal(t, restaurantID, result[0].RestaurantID)
	assert.Equal(t, restaurantID, result[1].RestaurantID)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_GetMenu_EmptyMenu(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	restaurantID := uuid.New()
	expectedMenu := []models.RestaurantMenuItem{}

	mockService.On("GetMenu", ctx, restaurantID).Return(expectedMenu, nil)

	result, err := uc.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenu, result)
	assert.Len(t, result, 0)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_GetMenu_ServiceError(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	restaurantID := uuid.New()
	expectedError := errors.New("failed to fetch menu")

	mockService.On("GetMenu", ctx, restaurantID).Return([]models.RestaurantMenuItem{}, expectedError)

	result, err := uc.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestRestaurantUseCase_GetMenu_InvalidRestaurantID(t *testing.T) {
	mockService := new(MockRestaurantService)
	uc := NewRestaurantUseCase(mockService)

	ctx := context.Background()
	restaurantID := uuid.New()
	expectedError := errors.New("restaurant not found")

	mockService.On("GetMenu", ctx, restaurantID).Return([]models.RestaurantMenuItem{}, expectedError)

	result, err := uc.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}
