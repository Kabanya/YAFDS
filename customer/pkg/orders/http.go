package orders

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"customer/pkg"

	"github.com/google/uuid"
)

type createRequest struct {
	CustomerID string `json:"customer_id"`
	CourierID  string `json:"courier_id"`
	Status     string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type courierResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewHandler(repo Repository) http.HandlerFunc {
	create := NewCreateHandler(repo)
	list := NewListHandler(repo)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r)
		case http.MethodGet:
			list(w, r)
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(errorResponse{Error: "method not allowed"})
		}
	}
}

func NewCreateHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := pkg.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req createRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			writeError(w, "customer_id must be UUID", http.StatusBadRequest)
			return
		}
		courierID, err := uuid.Parse(req.CourierID)
		if err != nil {
			writeError(w, "courier_id must be UUID", http.StatusBadRequest)
			return
		}

		created, err := repo.Create(r.Context(), Order{
			CustomerID: customerID,
			CourierID:  courierID,
			Status:     req.Status,
		})
		if err != nil {
			logger.Printf("orders: create failed: %v", err)
			switch {
			case errors.Is(err, ErrCustomerNotFound):
				writeError(w, "customer_id not found", http.StatusBadRequest)
			case errors.Is(err, ErrCourierNotFound):
				writeError(w, "courier_id not found", http.StatusBadRequest)
			default:
				writeError(w, "failed to create order", http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, created, http.StatusCreated)
	}
}

func NewListHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := pkg.Logger()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var filter Filter
		if v := r.URL.Query().Get("customer_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				writeError(w, "problem with customer_id (UUID)", http.StatusBadRequest)
				return
			}
			filter.CustomerID = &id
		}
		if v := r.URL.Query().Get("courier_id"); v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				writeError(w, "problem with courier_id (UUID)", http.StatusBadRequest)
				return
			}
			filter.CourierID = &id
		}
		if v := r.URL.Query().Get("status"); v != "" {
			filter.Status = v
		}

		orders, err := repo.List(r.Context(), filter)
		if err != nil {
			logger.Printf("orders: list failed: %v", err)
			writeError(w, "failed to fetch orders", http.StatusInternalServerError)
			return
		}

		writeJSON(w, orders, http.StatusOK)
	}
}

func NewCouriersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, _ := pkg.Logger()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.QueryContext(r.Context(), "SELECT emp_id, name FROM COURIERS WHERE is_active = TRUE")
		if err != nil {
			logger.Printf("orders: list couriers failed: %v", err)
			writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var couriers []courierResponse
		for rows.Next() {
			var c courierResponse
			if err := rows.Scan(&c.ID, &c.Name); err != nil {
				logger.Printf("orders: scan couriers failed: %v", err)
				writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
				return
			}
			couriers = append(couriers, c)
		}

		if err := rows.Err(); err != nil {
			logger.Printf("orders: iterate couriers failed: %v", err)
			writeError(w, "failed to fetch couriers", http.StatusInternalServerError)
			return
		}

		writeJSON(w, couriers, http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, errorResponse{Error: message}, status)
}
