package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kabanya/YAFDS/pkg/models"
	repositoryModels "github.com/Kabanya/YAFDS/pkg/repository/models"
	"github.com/google/uuid"
)

func TestPgRepo_CreateOrder(t *testing.T) {
	// Setup mocks
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	customersDB, customersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer customersDB.Close()

	couriersDB, couriersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer couriersDB.Close()

	repo := NewPostgresRepository(ordersDB, customersDB, couriersDB)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()
	// Let repository generate ID if not provided, or provide it to assert
	orderID := uuid.New()

	order := models.Order{
		ID:         orderID,
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     models.OrderStatusCustomerCreated,
	}

	// 1. ensureCustomerExists
	customersMock.ExpectQuery("SELECT 1 FROM customers WHERE emp_id = \\$1").
		WithArgs(customerID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 2. ensureCourierExists
	couriersMock.ExpectQuery("SELECT 1 FROM couriers WHERE emp_id = \\$1").
		WithArgs(courierID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 3. Insert Order
	// Note: created_at and updated_at are set inside the method, so we use AnyArg()
	ordersMock.ExpectExec("INSERT INTO ORDERS").
		WithArgs(orderID, customerID, courierID, sqlmock.AnyArg(), sqlmock.AnyArg(), string(models.OrderStatusCustomerCreated)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := repo.CreateOrder(ctx, order)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if created.ID != orderID {
		t.Errorf("expected order ID %v, got %v", orderID, created.ID)
	}

	if err := ordersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("ordersdb: %s", err)
	}
	if err := customersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("customersdb: %s", err)
	}
	if err := couriersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("couriersdb: %s", err)
	}
}

func TestPgRepo_CreateOrderWithItems(t *testing.T) {
	// Setup mocks
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	customersDB, customersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer customersDB.Close()

	couriersDB, couriersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer couriersDB.Close()

	repo := NewPostgresRepository(ordersDB, customersDB, couriersDB)

	ctx := context.Background()
	customerID := uuid.New()
	courierID := uuid.New()
	orderID := uuid.New()

	items := []repositoryModels.OrderItemInput{
		{
			RestaurantItemID: uuid.New(),
			Price:            10.0,
			Quantity:         2,
		},
		{
			RestaurantItemID: uuid.New(),
			Price:            5.0,
			Quantity:         1,
		},
	}

	order := models.Order{
		ID:         orderID,
		CustomerID: customerID,
		CourierID:  courierID,
		Status:     models.OrderStatusCustomerCreated,
	}

	// 1. ensureCustomerExists
	customersMock.ExpectQuery("SELECT 1 FROM customers WHERE emp_id = \\$1").
		WithArgs(customerID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 2. ensureCourierExists
	couriersMock.ExpectQuery("SELECT 1 FROM couriers WHERE emp_id = \\$1").
		WithArgs(courierID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 3. Begin TX
	ordersMock.ExpectBegin()

	// 4. Insert Order
	ordersMock.ExpectExec("INSERT INTO ORDERS").
		WithArgs(orderID, customerID, courierID, sqlmock.AnyArg(), sqlmock.AnyArg(), string(models.OrderStatusCustomerCreated)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 5. Insert Items
	for _, item := range items {
		ordersMock.ExpectExec("INSERT INTO ORDERS_ITEMS").
			WithArgs(sqlmock.AnyArg(), orderID, item.RestaurantItemID, item.Price, item.Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// 6. Commit TX
	ordersMock.ExpectCommit()

	created, err := repo.CreateOrderWithItems(ctx, order, items)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if created.ID != orderID {
		t.Errorf("expected order ID %v, got %v", orderID, created.ID)
	}

	if err := ordersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("ordersdb: %s", err)
	}
	if err := customersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("customersdb: %s", err)
	}
	if err := couriersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("couriersdb: %s", err)
	}
}

func TestPgRepo_GetOrder(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)

	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()
	courierID := uuid.New()
	now := time.Now().UTC()
	status := models.OrderStatusCustomerCreated

	// QueryRowContext
	ordersMock.ExpectQuery("SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"emp_id", "customer_id", "courier_id", "created_at", "updated_at", "status"}).
			AddRow(orderID, customerID, courierID, now, now, string(status)))

	o, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if o.ID != orderID {
		t.Errorf("expected orderID %v, got %v", orderID, o.ID)
	}
	if o.Status != status {
		t.Errorf("expected status %v, got %v", status, o.Status)
	}

	// Test Not Found
	ordersMock.ExpectQuery("SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetOrder(ctx, orderID)
	if err != ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestPgRepo_ListOrders(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)
	ctx := context.Background()

	// 1. List with no filter
	orderID1 := uuid.New()
	orderID2 := uuid.New()
	now := time.Now().UTC()

	ordersMock.ExpectQuery("SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS ORDER BY created_at DESC").
		WillReturnRows(sqlmock.NewRows([]string{"emp_id", "customer_id", "courier_id", "created_at", "updated_at", "status"}).
			AddRow(orderID1, uuid.New(), uuid.New(), now, now, "STATUS1").
			AddRow(orderID2, uuid.New(), uuid.New(), now, now, "STATUS2"))

	list, err := repo.ListOrders(ctx, repositoryModels.Filter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 orders, got %d", len(list))
	}

	// 2. List with filter
	custID := uuid.New()
	filter := repositoryModels.Filter{CustomerID: &custID}

	ordersMock.ExpectQuery("SELECT emp_id, customer_id, courier_id, created_at, updated_at, status FROM ORDERS WHERE customer_id = \\$1 ORDER BY created_at DESC").
		WithArgs(custID).
		WillReturnRows(sqlmock.NewRows([]string{"emp_id", "customer_id", "courier_id", "created_at", "updated_at", "status"}).
			AddRow(orderID1, custID, uuid.New(), now, now, "STATUS1"))

	listFiltered, err := repo.ListOrders(ctx, filter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(listFiltered) != 1 {
		t.Errorf("expected 1 order, got %d", len(listFiltered))
	}
}

func TestPgRepo_AcceptOrder(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	customersDB, customersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer customersDB.Close()

	couriersDB, couriersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer couriersDB.Close()

	repo := NewPostgresRepository(ordersDB, customersDB, couriersDB)
	ctx := context.Background()

	input := repositoryModels.AcceptInput{
		OrderID:    uuid.New(),
		CustomerID: uuid.New(),
		CourierID:  uuid.New(),
		Status:     models.OrderStatusKitchenAccepted,
		Items:      nil,
	}

	// 1. Check customer
	customersMock.ExpectQuery("SELECT 1 FROM customers").
		WithArgs(input.CustomerID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 2. Check courier
	couriersMock.ExpectQuery("SELECT 1 FROM couriers").
		WithArgs(input.CourierID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// 3. Begin TX
	ordersMock.ExpectBegin()

	// 4. Check status
	ordersMock.ExpectQuery("SELECT status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(input.OrderID).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("CUSTOMER_CREATED"))

	// 5. Upsert Order
	ordersMock.ExpectExec("INSERT INTO ORDERS").
		WithArgs(input.OrderID, input.CustomerID, input.CourierID, sqlmock.AnyArg(), sqlmock.AnyArg(), string(input.Status)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 6. Count items
	ordersMock.ExpectQuery("SELECT COUNT\\(1\\) FROM ORDERS_ITEMS WHERE order_id = \\$1").
		WithArgs(input.OrderID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// 7. Commit
	ordersMock.ExpectCommit()

	res, err := repo.AcceptOrder(ctx, input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if res.Status != string(input.Status) {
		t.Errorf("expected status %v, got %v", input.Status, res.Status)
	}
}

func TestPgRepo_GetOrderStatus(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)
	ctx := context.Background()
	orderID := uuid.New()
	status := models.OrderStatusKitchenAccepted

	// Test successful retrieval
	ordersMock.ExpectQuery("SELECT status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow(string(status)))

	gotStatus, err := repo.GetOrderStatus(ctx, orderID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gotStatus != status {
		t.Errorf("expected status %v, got %v", status, gotStatus)
	}

	// Test order not found
	ordersMock.ExpectQuery("SELECT status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetOrderStatus(ctx, orderID)
	if err != ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}

	// Test database error
	expectedErr := sql.ErrConnDone
	ordersMock.ExpectQuery("SELECT status FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnError(expectedErr)

	_, err = repo.GetOrderStatus(ctx, orderID)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if err := ordersMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgRepo_UpdateOrderStatus(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)
	ctx := context.Background()
	orderID := uuid.New()
	newStatus := models.OrderStatusDeliveryDelivering

	// Success
	ordersMock.ExpectExec("UPDATE ORDERS SET status = \\$1, updated_at = \\$2 WHERE emp_id = \\$3").
		WithArgs(string(newStatus), sqlmock.AnyArg(), orderID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Not Found (rows affected 0)
	ordersMock.ExpectExec("UPDATE ORDERS SET status = \\$1, updated_at = \\$2 WHERE emp_id = \\$3").
		WithArgs(string(newStatus), sqlmock.AnyArg(), orderID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateOrderStatus(ctx, orderID, newStatus)
	if err != ErrOrderNotFound {
		t.Errorf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestPgRepo_CalculateOrderTotal(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)
	ctx := context.Background()
	orderID := uuid.New()

	// 1. Result found
	ordersMock.ExpectQuery("SELECT SUM\\(price \\* quantity\\) FROM ORDERS_ITEMS WHERE order_id = \\$1").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(150.50))

	total, err := repo.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if total != 150.50 {
		t.Errorf("expected 150.50, got %f", total)
	}

	// 2. No items (NULL sum)
	ordersMock.ExpectQuery("SELECT SUM\\(price \\* quantity\\) FROM ORDERS_ITEMS WHERE order_id = \\$1").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(nil))

	total, err = repo.CalculateOrderTotal(ctx, orderID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Errorf("expected 0, got %f", total)
	}
}

func TestPgRepo_AddItemIntoOrder(t *testing.T) {
	ordersDB, ordersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer ordersDB.Close()

	repo := NewPostgresRepository(ordersDB, nil, nil)
	ctx := context.Background()
	orderID := uuid.New()
	restItemID := uuid.New()

	item := repositoryModels.OrderItemInput{
		RestaurantItemID: restItemID,
		Price:            10.0,
		Quantity:         2,
	}

	ordersMock.ExpectBegin()
	// Check exists
	ordersMock.ExpectQuery("SELECT 1 FROM ORDERS WHERE emp_id = \\$1").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	// Insert item
	ordersMock.ExpectExec("INSERT INTO ORDERS_ITEMS").
		WithArgs(sqlmock.AnyArg(), orderID, restItemID, 10.0, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Update order updated_at
	ordersMock.ExpectExec("UPDATE ORDERS SET updated_at = \\$1 WHERE emp_id = \\$2").
		WithArgs(sqlmock.AnyArg(), orderID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ordersMock.ExpectCommit()

	err = repo.AddItemIntoOrder(ctx, orderID, item)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPgRepo_GetCustomerWalletAddress(t *testing.T) {
	customersDB, customersMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer customersDB.Close()

	repo := NewPostgresRepository(nil, customersDB, nil)
	ctx := context.Background()
	customerID := uuid.New()
	wallet := "0x123456789"

	customersMock.ExpectQuery("SELECT wallet_address FROM CUSTOMERS WHERE emp_id = \\$1").
		WithArgs(customerID).
		WillReturnRows(sqlmock.NewRows([]string{"wallet_address"}).AddRow(wallet))

	got, err := repo.GetCustomerWalletAddress(ctx, customerID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != wallet {
		t.Errorf("expected wallet %v, got %v", wallet, got)
	}

	// Not found
	customersMock.ExpectQuery("SELECT wallet_address FROM CUSTOMERS WHERE emp_id = \\$1").
		WithArgs(customerID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetCustomerWalletAddress(ctx, customerID)
	if err != ErrCustomerNotFound {
		t.Errorf("expected ErrCustomerNotFound, got %v", err)
	}
}
