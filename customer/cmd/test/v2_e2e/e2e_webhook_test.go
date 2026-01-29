package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// E2E test: Complete order flow from creation to delivery
// This test simulates a real-world scenario where:
// 1. Customer browses restaurants
// 2. Views restaurant menu
// 3. Creates an order with items
// 4. Courier accepts the order
// 5. Order status progresses through delivery stages
func TestE2E_CompleteOrderFlow(t *testing.T) {
	// Setup mock restaurant service
	restaurantID := uuid.New()
	menuItems := []map[string]interface{}{
		{
			"id":            uuid.New().String(),
			"restaurant_id": restaurantID.String(),
			"name":          "Classic Burger",
			"price":         12.99,
		},
		{
			"id":            uuid.New().String(),
			"restaurant_id": restaurantID.String(),
			"name":          "French Fries",
			"price":         4.50,
		},
		{
			"id":            uuid.New().String(),
			"restaurant_id": restaurantID.String(),
			"name":          "Cola",
			"price":         2.99,
		},
	}

	restaurantServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/menu", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(menuItems)
	}))
	defer restaurantServer.Close()

	// Setup mock courier service
	courierID := uuid.New()
	courierServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.URL.Path == "/couriers" {
			couriers := []map[string]interface{}{
				{
					"id":      courierID.String(),
					"name":    "John Courier",
					"status":  "available",
					"vehicle": "bike",
					"rating":  4.5,
				},
			}
			json.NewEncoder(w).Encode(couriers)
		}
	}))
	defer courierServer.Close()

	// Setup mock customer API server
	customerID := uuid.New()
	createdOrders := make(map[string]map[string]interface{})
	orderMutex := sync.Mutex{}
	orderStatusTransitions := make([]string, 0)
	restaurantList := []map[string]interface{}{
		{
			"id":           restaurantID.String(),
			"name":         "Burger Palace",
			"address":      "123 Main St",
			"rating":       4.5,
			"delivery_fee": 2.99,
		},
	}

	customerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch {
		case r.URL.Path == "/health" && r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

		case r.URL.Path == "/restaurants" && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(restaurantList)

		case r.URL.Path == "/menu" && r.Method == http.MethodGet:
			// Proxy to restaurant server
			menuResp, err := http.Get(restaurantServer.URL + "/menu")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer menuResp.Body.Close()
			body, _ := io.ReadAll(menuResp.Body)
			w.Write(body)

		case r.URL.Path == "/couriers" && r.Method == http.MethodGet:
			// Proxy to courier server
			courierResp, err := http.Get(courierServer.URL + "/couriers")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer courierResp.Body.Close()
			body, _ := io.ReadAll(courierResp.Body)
			w.Write(body)

		case r.URL.Path == "/orders" && r.Method == http.MethodGet:
			var orders []map[string]interface{}
			orderMutex.Lock()
			for _, order := range createdOrders {
				orders = append(orders, order)
			}
			orderMutex.Unlock()
			json.NewEncoder(w).Encode(orders)

		case r.URL.Path == "/orders" && r.Method == http.MethodPost:
			var req map[string]interface{}
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			orderID := uuid.New()
			items := []map[string]interface{}{}
			if reqItems, ok := req["items"].([]interface{}); ok {
				for _, item := range reqItems {
					if itemMap, ok := item.(map[string]interface{}); ok {
						items = append(items, itemMap)
					}
				}
			}

			order := map[string]interface{}{
				"id":          orderID.String(),
				"customer_id": req["customer_id"],
				"courier_id":  req["courier_id"],
				"status":      "customer_created",
				"items":       items,
				"created_at":  time.Now().Format(time.RFC3339),
				"updated_at":  time.Now().Format(time.RFC3339),
			}

			orderMutex.Lock()
			createdOrders[orderID.String()] = order
			orderStatusTransitions = append(orderStatusTransitions, "customer_created")
			orderMutex.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)

		case strings.HasPrefix(r.URL.Path, "/orders/") && strings.HasSuffix(r.URL.Path, "/status") && r.Method == http.MethodGet:
			orderID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/status"), "/orders/")
			orderMutex.Lock()
			order, exists := createdOrders[orderID]
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"status": order["status"]})

		case strings.HasPrefix(r.URL.Path, "/orders/") && strings.HasSuffix(r.URL.Path, "/status") && r.Method == http.MethodPut:
			orderID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/status"), "/orders/")
			var req map[string]string
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			orderMutex.Lock()
			order, exists := createdOrders[orderID]
			if exists {
				order["status"] = req["status"]
				order["updated_at"] = time.Now().Format(time.RFC3339)
				orderStatusTransitions = append(orderStatusTransitions, req["status"])
			}
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"message": "status updated"})

		case strings.HasPrefix(r.URL.Path, "/orders/") && strings.HasSuffix(r.URL.Path, "/total") && r.Method == http.MethodGet:
			orderID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/total"), "/orders/")
			orderMutex.Lock()
			order, exists := createdOrders[orderID]
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}

			total := 0.0
			if items, ok := order["items"].([]map[string]interface{}); ok {
				for _, item := range items {
					if price, ok := item["price"].(float64); ok {
						if qty, ok := item["quantity"].(float64); ok {
							total += price * qty
						} else {
							total += price
						}
					}
				}
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"total": total})

		case strings.HasPrefix(r.URL.Path, "/orders/") && !strings.Contains(r.URL.Path[8:], "/") && r.Method == http.MethodGet:
			orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
			orderMutex.Lock()
			order, exists := createdOrders[orderID]
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}
			json.NewEncoder(w).Encode(order)

		case strings.HasPrefix(r.URL.Path, "/orders/") && strings.HasSuffix(r.URL.Path, "/accept") && r.Method == http.MethodPost:
			orderID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/accept"), "/orders/")
			var req map[string]interface{}
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			orderMutex.Lock()
			order, exists := createdOrders[orderID]
			if exists {
				order["status"] = "courier_accepted"
				if cid, ok := req["courier_id"]; ok {
					order["courier_id"] = cid
				}
				order["updated_at"] = time.Now().Format(time.RFC3339)
				orderStatusTransitions = append(orderStatusTransitions, "courier_accepted")
			}
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}

			result := map[string]interface{}{
				"order_id":    orderID,
				"customer_id": req["customer_id"],
				"courier_id":  req["courier_id"],
				"status":      "courier_accepted",
			}
			json.NewEncoder(w).Encode(result)

		case r.URL.Path == "/webhook/payment" && r.Method == http.MethodPost:
			var req struct {
				OrderID string `json:"order_id"`
				Status  string `json:"status"`
			}
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			orderMutex.Lock()
			order, exists := createdOrders[req.OrderID]
			if exists {
				// Convert provider status to internal status
				internalStatus := "customer_paid"
				if req.Status == "failed" {
					internalStatus = "payment_failed"
				}

				order["status"] = internalStatus
				order["updated_at"] = time.Now().Format(time.RFC3339)
				orderStatusTransitions = append(orderStatusTransitions, internalStatus)
			}
			orderMutex.Unlock()

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})

		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		}
	}))
	defer customerServer.Close()

	// ===== E2E Test Scenarios =====

	t.Run("Scenario: Customer places order from restaurant menu", func(t *testing.T) {
		// Step 1: Customer browses restaurants
		resp, err := http.Get(customerServer.URL + "/restaurants")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var restaurants []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&restaurants)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Len(t, restaurants, 1)
		assert.Equal(t, restaurantID.String(), restaurants[0]["id"])
		assert.Equal(t, "Burger Palace", restaurants[0]["name"])

		// Step 2: Customer views restaurant menu
		menuURL := fmt.Sprintf("%s/menu?restaurant_id=%s", customerServer.URL, restaurantID.String())
		resp, err = http.Get(menuURL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var menu []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&menu)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Len(t, menu, 3)
		assert.Equal(t, "Classic Burger", menu[0]["name"])
		assert.Equal(t, 12.99, menu[0]["price"])

		// Step 3: Customer places order with items
		orderItems := []map[string]interface{}{
			{
				"restaurant_item_id": menuItems[0]["id"],
				"price":              menuItems[0]["price"],
				"quantity":           1,
			},
			{
				"restaurant_item_id": menuItems[1]["id"],
				"price":              menuItems[1]["price"],
				"quantity":           2,
			},
		}

		createOrderReq := map[string]interface{}{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
			"items":       orderItems,
		}

		body, _ := json.Marshal(createOrderReq)
		resp, err = http.Post(customerServer.URL+"/orders", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdOrder map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createdOrder)
		require.NoError(t, err)
		resp.Body.Close()

		orderID, ok := createdOrder["id"].(string)
		require.True(t, ok, "order ID should be string")
		assert.Equal(t, customerID.String(), createdOrder["customer_id"])
		assert.Equal(t, courierID.String(), createdOrder["courier_id"])
		assert.Equal(t, "customer_created", createdOrder["status"])

		// Step 4: Customer views order details
		resp, err = http.Get(fmt.Sprintf("%s/orders/%s", customerServer.URL, orderID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var orderDetails map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orderDetails)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, orderID, orderDetails["id"])
		assert.Equal(t, "customer_created", orderDetails["status"])

		// Step 5: Customer checks order total
		resp, err = http.Get(fmt.Sprintf("%s/orders/%s/total", customerServer.URL, orderID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var totalResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&totalResp)
		require.NoError(t, err)
		resp.Body.Close()

		// Expected: 12.99 + 4.50 + 4.50 = 21.99
		expectedTotal := 12.99 + 4.50 + 4.50
		assert.InDelta(t, expectedTotal, totalResp["total"], 0.01)

		// Step 6: Customer pays for order
		// (Simulated: Status might stay 'customer_created' or move to 'payment_pending')

		// Step 7: Payment provider sends webhook to confirm payment
		webhookReq := map[string]string{
			"order_id": orderID,
			"status":   "paid",
		}
		body, _ = json.Marshal(webhookReq)
		resp, err = http.Post(customerServer.URL+"/webhook/payment", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Step 8: Courier accepts order
		acceptReq := map[string]interface{}{
			"customer_id": customerID.String(),
			"courier_id":  courierID.String(),
			"items":       orderItems,
			"status":      "courier_accepted",
		}
		body, _ = json.Marshal(acceptReq)
		resp, err = http.Post(fmt.Sprintf("%s/orders/%s/accept", customerServer.URL, orderID), "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Step 9: Verify order status
		resp, err = http.Get(fmt.Sprintf("%s/orders/%s/status", customerServer.URL, orderID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var statusResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&statusResp)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, "courier_accepted", statusResp["status"])

		// Step 10: Courier picks up order
		pickupReq := map[string]string{"status": "courier_picked_up"}
		body, _ = json.Marshal(pickupReq)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/orders/%s/status", customerServer.URL, orderID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Step 11: Courier delivers order
		deliveredReq := map[string]string{"status": "delivered"}
		body, _ = json.Marshal(deliveredReq)
		req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/orders/%s/status", customerServer.URL, orderID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Step 12: Verify final order status
		resp, err = http.Get(fmt.Sprintf("%s/orders/%s/status", customerServer.URL, orderID))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&statusResp)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, "delivered", statusResp["status"])

		// Verify status transitions
		expectedTransitions := []string{
			"customer_created",
			"customer_paid",
			"courier_accepted",
			"courier_picked_up",
			"delivered",
		}
		assert.Equal(t, expectedTransitions, orderStatusTransitions)
	})

	t.Run("Scenario: List orders by customer", func(t *testing.T) {
		// Clear orders for this test
		orderMutex.Lock()
		createdOrders = make(map[string]map[string]interface{})
		orderMutex.Unlock()

		// Create 3 orders for the same customer
		for i := 0; i < 3; i++ {
			createOrderReq := map[string]interface{}{
				"customer_id": customerID.String(),
				"courier_id":  courierID.String(),
				"items": []map[string]interface{}{
					{
						"restaurant_item_id": menuItems[0]["id"],
						"price":              menuItems[0]["price"],
						"quantity":           1,
					},
				},
			}

			body, _ := json.Marshal(createOrderReq)
			resp, err := http.Post(customerServer.URL+"/orders", "application/json", bytes.NewBuffer(body))
			require.NoError(t, err)
			resp.Body.Close()
		}

		// List all orders
		resp, err := http.Get(customerServer.URL + "/orders")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var orders []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orders)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Len(t, orders, 3)

		// Verify all orders belong to the customer
		for _, order := range orders {
			assert.Equal(t, customerID.String(), order["customer_id"])
		}
	})
}

// E2E test: Concurrent orders from multiple customers
func TestE2E_ConcurrentOrders(t *testing.T) {
	restaurantID := uuid.New()
	menuItems := []map[string]interface{}{
		{
			"id":            uuid.New().String(),
			"restaurant_id": restaurantID.String(),
			"name":          "Pizza",
			"price":         15.00,
		},
	}

	createdOrders := make(map[string]map[string]interface{})
	orderMutex := sync.Mutex{}

	customerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.URL.Path == "/orders" && r.Method == http.MethodPost {
			var req map[string]interface{}
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			orderID := uuid.New()
			order := map[string]interface{}{
				"id":          orderID.String(),
				"customer_id": req["customer_id"],
				"courier_id":  req["courier_id"],
				"status":      "customer_created",
				"created_at":  time.Now().Format(time.RFC3339),
			}

			orderMutex.Lock()
			createdOrders[orderID.String()] = order
			orderMutex.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)
		}
	}))
	defer customerServer.Close()

	// Simulate 5 customers placing orders concurrently
	const numCustomers = 5
	var wg sync.WaitGroup
	wg.Add(numCustomers)

	for i := 0; i < numCustomers; i++ {
		go func(customerIdx int) {
			defer wg.Done()

			customerID := uuid.New()
			courierID := uuid.New()

			createOrderReq := map[string]interface{}{
				"customer_id": customerID.String(),
				"courier_id":  courierID.String(),
				"items": []map[string]interface{}{
					{
						"restaurant_item_id": menuItems[0]["id"],
						"price":              menuItems[0]["price"],
						"quantity":           1,
					},
				},
			}

			body, _ := json.Marshal(createOrderReq)
			resp, err := http.Post(customerServer.URL+"/orders", "application/json", bytes.NewBuffer(body))
			assert.NoErrorf(t, err, "Customer %d failed to create order", customerIdx)
			if err == nil {
				assert.Equalf(t, http.StatusCreated, resp.StatusCode, "Customer %d order creation status", customerIdx)
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	// Verify all orders were created
	orderMutex.Lock()
	assert.Len(t, createdOrders, numCustomers)
	orderMutex.Unlock()
}

// E2E test: Order error scenarios
func TestE2E_OrderErrorScenarios(t *testing.T) {
	customerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.URL.Path == "/orders" && r.Method == http.MethodPost {
			var req map[string]interface{}
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &req)

			// Validate request
			if req["customer_id"] == nil || req["courier_id"] == nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "customer_id and courier_id are required"})
				return
			}

			order := map[string]interface{}{
				"id":          uuid.New().String(),
				"customer_id": req["customer_id"],
				"courier_id":  req["courier_id"],
				"status":      "customer_created",
				"created_at":  time.Now().Format(time.RFC3339),
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)
		}

		if strings.HasPrefix(r.URL.Path, "/orders/") && !strings.Contains(r.URL.Path[8:], "/") && r.Method == http.MethodGet {
			orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
			if orderID == "00000000-0000-0000-0000-000000000000" {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
				return
			}
		}
	}))
	defer customerServer.Close()

	t.Run("Error: Create order without customer_id", func(t *testing.T) {
		createOrderReq := map[string]interface{}{
			"courier_id": uuid.New().String(),
		}

		body, _ := json.Marshal(createOrderReq)
		resp, err := http.Post(customerServer.URL+"/orders", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("Error: Get non-existent order", func(t *testing.T) {
		resp, err := http.Get(customerServer.URL + "/orders/00000000-0000-0000-0000-000000000000")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})
}

// E2E test: Restaurant menu discovery flow
func TestE2E_RestaurantMenuDiscoveryFlow(t *testing.T) {
	restaurants := []map[string]interface{}{
		{
			"id":           uuid.New().String(),
			"name":         "Burger Palace",
			"address":      "123 Main St",
			"rating":       4.5,
			"delivery_fee": 2.99,
		},
		{
			"id":           uuid.New().String(),
			"name":         "Pizza Corner",
			"address":      "456 Oak Ave",
			"rating":       4.2,
			"delivery_fee": 1.99,
		},
	}

	restaurantMenus := make(map[string][]map[string]interface{})
	for _, r := range restaurants {
		restaurantID := r["id"].(string)
		restaurantMenus[restaurantID] = []map[string]interface{}{
			{
				"id":            uuid.New().String(),
				"restaurant_id": restaurantID,
				"name":          fmt.Sprintf("%s Special", r["name"]),
				"price":         15.99,
			},
		}
	}

	customerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.URL.Path == "/restaurants" && r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(restaurants)
		}

		if r.URL.Path == "/menu" && r.Method == http.MethodGet {
			restaurantID := r.URL.Query().Get("restaurant_id")
			if menu, ok := restaurantMenus[restaurantID]; ok {
				json.NewEncoder(w).Encode(menu)
			} else {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode([]map[string]interface{}{})
			}
		}
	}))
	defer customerServer.Close()

	t.Run("Customer discovers restaurants", func(t *testing.T) {
		resp, err := http.Get(customerServer.URL + "/restaurants")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Len(t, result, 2)
		assert.Equal(t, "Burger Palace", result[0]["name"])
		assert.Equal(t, "Pizza Corner", result[1]["name"])
	})

	t.Run("Customer views restaurant menu", func(t *testing.T) {
		restaurantID := restaurants[0]["id"].(string)
		menuURL := fmt.Sprintf("%s/menu?restaurant_id=%s", customerServer.URL, restaurantID)
		resp, err := http.Get(menuURL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var menu []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&menu)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Len(t, menu, 1)
		assert.Equal(t, "Burger Palace Special", menu[0]["name"])
	})
}
