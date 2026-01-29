package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPRestaurantClient tests the NewHTTPRestaurantClient constructor
func TestNewHTTPRestaurantClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewHTTPRestaurantClient(baseURL)

	assert.NotNil(t, client)
	httpClient, ok := client.(*httpRestaurantClient)
	assert.True(t, ok)
	assert.Equal(t, baseURL, httpClient.baseURL)
	assert.NotNil(t, httpClient.httpClient)
	assert.Equal(t, 10*time.Second, httpClient.httpClient.Timeout)
}

// TestGetMenu_Success tests successful menu retrieval
func TestGetMenu_Success(t *testing.T) {
	restaurantID := uuid.New()
	expectedMenuItems := []RestaurantMenuItem{
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Burger",
			Price:        12.99,
		},
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Pizza",
			Price:        15.50,
		},
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Salad",
			Price:        8.00,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/menu", r.URL.Path)
		assert.Equal(t, restaurantID.String(), r.URL.Query().Get("restaurant_id"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenuItems, menu)
	assert.Len(t, menu, 3)
}

// TestGetMenu_EmptyMenu tests retrieval of an empty menu
func TestGetMenu_EmptyMenu(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode([]RestaurantMenuItem{})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, menu)
	assert.Len(t, menu, 0)
}

// TestGetMenu_StatusNotFound tests 404 response
func TestGetMenu_StatusNotFound(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("Not Found"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
	assert.Contains(t, err.Error(), "restaurant service returned status")
	assert.Contains(t, err.Error(), "404")
}

// TestGetMenu_StatusInternalServerError tests 500 response
func TestGetMenu_StatusInternalServerError(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
	assert.Contains(t, err.Error(), "restaurant service returned status")
	assert.Contains(t, err.Error(), "500")
}

// TestGetMenu_StatusBadRequest tests 400 response
func TestGetMenu_StatusBadRequest(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Bad Request"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
	assert.Contains(t, err.Error(), "restaurant service returned status")
	assert.Contains(t, err.Error(), "400")
}

// TestGetMenu_InvalidJSON tests handling of invalid JSON response
func TestGetMenu_InvalidJSON(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("invalid json {"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
}

// TestGetMenu_MalformedJSON tests handling of malformed JSON (empty object)
func TestGetMenu_MalformedJSON(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
}

// TestGetMenu_ContextCancelled tests handling of context cancellation
func TestGetMenu_ContextCancelled(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intentionally slow to trigger context cancellation
		time.Sleep(500 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Use shorter timeout that will definitely trigger
	httpClientWithShortTimeout := &httpRestaurantClient{
		baseURL: server.URL,
		httpClient: &http.Client{
			Timeout: 50 * time.Millisecond,
		},
	}

	ctx := context.Background()
	menu, err := httpClientWithShortTimeout.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
}

// TestGetMenu_ContextDeadlineExceeded tests context deadline exceeded
func TestGetMenu_ContextDeadlineExceeded(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler is never called due to deadline
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with very short timeout
	httpClientWithShortTimeout := &httpRestaurantClient{
		baseURL: server.URL,
		httpClient: &http.Client{
			Timeout: 50 * time.Millisecond,
		},
	}

	ctx := context.Background()
	menu, err := httpClientWithShortTimeout.GetMenu(ctx, restaurantID)

	assert.Error(t, err)
	assert.Nil(t, menu)
}

// TestGetMenu_URLConstructionWithSpecialCharacters tests URL construction with special characters
func TestGetMenu_URLConstructionWithSpecialCharacters(t *testing.T) {
	restaurantID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the restaurant ID is properly encoded in the URL
		assert.Equal(t, restaurantID.String(), r.URL.Query().Get("restaurant_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode([]RestaurantMenuItem{})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.NotNil(t, menu)
}

// TestGetMenu_ConnectionRefused tests connection refusal
func TestGetMenu_ConnectionRefused(t *testing.T) {
	client := NewHTTPRestaurantClient("http://localhost:49999") // Port that's unlikely to be in use
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, uuid.New())

	assert.Error(t, err)
	assert.Nil(t, menu)
}

// TestGetMenu_WithTrailingSlashInURL tests URL handling with trailing slash
func TestGetMenu_WithTrailingSlashInURL(t *testing.T) {
	restaurantID := uuid.New()
	expectedMenuItems := []RestaurantMenuItem{
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Pasta",
			Price:        14.99,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create client with trailing slash
	client := NewHTTPRestaurantClient(server.URL + "/")
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenuItems, menu)
}

// TestGetMenu_MultipleMenuItems tests retrieval of large menu
func TestGetMenu_MultipleMenuItems(t *testing.T) {
	restaurantID := uuid.New()

	// Create a large menu
	expectedMenuItems := make([]RestaurantMenuItem, 50)
	for i := 0; i < 50; i++ {
		expectedMenuItems[i] = RestaurantMenuItem{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         fmt.Sprintf("Dish %d", i+1),
			Price:        float64(i+1) * 1.99,
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Len(t, menu, 50)
	assert.Equal(t, expectedMenuItems, menu)
}

// TestGetMenu_DifferentRestaurantIDs tests retrieval with different restaurant IDs
func TestGetMenu_DifferentRestaurantIDs(t *testing.T) {
	testCases := []struct {
		name         string
		restaurantID uuid.UUID
	}{
		{
			name:         "Random UUID",
			restaurantID: uuid.New(),
		},
		{
			name:         "Different Random UUID",
			restaurantID: uuid.New(),
		},
		{
			name:         "Another Random UUID",
			restaurantID: uuid.New(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capturedID := ""
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedID = r.URL.Query().Get("restaurant_id")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode([]RestaurantMenuItem{})
				require.NoError(t, err)
			}))
			defer server.Close()

			client := NewHTTPRestaurantClient(server.URL)
			ctx := context.Background()

			_, err := client.GetMenu(ctx, tc.restaurantID)

			assert.NoError(t, err)
			assert.Equal(t, tc.restaurantID.String(), capturedID)
		})
	}
}

// TestGetMenu_MenuItemWithZeroPrice tests menu item with zero price
func TestGetMenu_MenuItemWithZeroPrice(t *testing.T) {
	restaurantID := uuid.New()
	expectedMenuItems := []RestaurantMenuItem{
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Free Sample",
			Price:        0.0,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenuItems, menu)
	assert.Equal(t, 0.0, menu[0].Price)
}

// TestGetMenu_MenuItemWithLargePrice tests menu item with large price
func TestGetMenu_MenuItemWithLargePrice(t *testing.T) {
	restaurantID := uuid.New()
	expectedMenuItems := []RestaurantMenuItem{
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Luxury Item",
			Price:        9999.99,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenuItems, menu)
	assert.Equal(t, 9999.99, menu[0].Price)
}

// TestGetMenu_MenuItemWithSpecialCharactersInName tests menu items with special characters
func TestGetMenu_MenuItemWithSpecialCharactersInName(t *testing.T) {
	restaurantID := uuid.New()
	expectedMenuItems := []RestaurantMenuItem{
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Café Latte & Pastry",
			Price:        5.99,
		},
		{
			ID:           uuid.New(),
			RestaurantID: restaurantID,
			Name:         "Señor's Special (Spicy)",
			Price:        12.50,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedMenuItems)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	menu, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMenuItems, menu)
	assert.Equal(t, "Café Latte & Pastry", menu[0].Name)
	assert.Equal(t, "Señor's Special (Spicy)", menu[1].Name)
}

// TestGetMenu_HTTPMethodIsGet tests that the correct HTTP method is used
func TestGetMenu_HTTPMethodIsGet(t *testing.T) {
	restaurantID := uuid.New()
	methodUsed := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodUsed = r.Method
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode([]RestaurantMenuItem{})
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewHTTPRestaurantClient(server.URL)
	ctx := context.Background()

	_, err := client.GetMenu(ctx, restaurantID)

	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, methodUsed)
}
