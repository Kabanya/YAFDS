package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RestaurantMenuItem struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
}

type RestaurantClient interface {
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error)
}

type httpRestaurantClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPRestaurantClient(baseURL string) RestaurantClient {
	return &httpRestaurantClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *httpRestaurantClient) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error) {
	url := fmt.Sprintf("%s/menu?restaurant_id=%s", c.baseURL, restaurantID.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("restaurant service returned status: %s", resp.Status)
	}

	var menu []RestaurantMenuItem
	if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
		return nil, err
	}

	return menu, nil
}
