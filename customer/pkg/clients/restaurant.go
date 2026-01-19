package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RestaurantMenuItem struct {
	OrderItemID  uuid.UUID `json:"order_item_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	Description  string    `json:"description"`
}

type RestaurantMenuClient interface {
	GetMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error)
}

type HTTPRestaurantClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPRestaurantClient(baseURL string) *HTTPRestaurantClient {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		trimmed = "http://localhost:8092"
	}
	return &HTTPRestaurantClient{
		baseURL: trimmed,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPRestaurantClient) GetMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error) {
	endpoint, err := url.Parse(c.baseURL + "/menu/show")
	if err != nil {
		return nil, err
	}
	query := endpoint.Query()
	query.Set("restaurant_id", restaurantID.String())
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		if errBody.Error == "" {
			errBody.Error = resp.Status
		}
		return nil, fmt.Errorf("restaurant menu request failed: %s", errBody.Error)
	}

	var items []RestaurantMenuItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
