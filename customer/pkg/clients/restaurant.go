package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	mu         sync.RWMutex
	cache      map[uuid.UUID]cachedMenu
}

type cachedMenu struct {
	items     []RestaurantMenuItem
	fetchedAt time.Time
}

const menuCacheTTL = 10 * time.Minute

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
		cache: make(map[uuid.UUID]cachedMenu),
	}
}

func (c *HTTPRestaurantClient) GetMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error) {
	now := time.Now().UTC()
	if items, ok := c.getCachedMenu(restaurantID, now); ok {
		return items, nil
	}

	items, err := c.fetchMenuItems(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	c.setCachedMenu(restaurantID, items, now)
	return items, nil
}

func (c *HTTPRestaurantClient) getCachedMenu(restaurantID uuid.UUID, now time.Time) ([]RestaurantMenuItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[restaurantID]
	if !ok {
		return nil, false
	}
	if now.Sub(entry.fetchedAt) > menuCacheTTL {
		return nil, false
	}
	return entry.items, true
}

func (c *HTTPRestaurantClient) setCachedMenu(restaurantID uuid.UUID, items []RestaurantMenuItem, fetchedAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = make(map[uuid.UUID]cachedMenu)
	}
	c.cache[restaurantID] = cachedMenu{items: items, fetchedAt: fetchedAt}
}

func (c *HTTPRestaurantClient) fetchMenuItems(ctx context.Context, restaurantID uuid.UUID) ([]RestaurantMenuItem, error) {
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
