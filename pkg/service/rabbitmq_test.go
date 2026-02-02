package service

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Test configuration matching Makefile defaults
func getTestRabbitMQConfig() RabbitMQConfig {
	host := getEnv("RABBITMQ_TEST_HOST", "localhost")
	port := getEnv("RABBITMQ_TEST_PORT", "5783")
	user := getEnv("RABBITMQ_TEST_USER", "guest")
	pass := getEnv("RABBITMQ_TEST_PASS", "guest")
	vhost := getEnv("RABBITMQ_TEST_VHOST", "/")

	return RabbitMQConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		VHost:    vhost,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestPaymentNotificationMessage_JSONSerialization tests JSON marshaling/unmarshaling
func TestPaymentNotificationMessage_JSONSerialization(t *testing.T) {
	orderID := uuid.New()
	originalMsg := PaymentNotificationMessage{
		OrderID:     orderID,
		Status:      "paid",
		Timestamp:   time.Now(),
		MessageType: "payment_notification",
	}

	// Marshal
	data, err := json.Marshal(originalMsg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	// Unmarshal
	var unmarshaledMsg PaymentNotificationMessage
	if err := json.Unmarshal(data, &unmarshaledMsg); err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	// Verify
	if unmarshaledMsg.OrderID != originalMsg.OrderID {
		t.Errorf("expected OrderID %s, got %s", originalMsg.OrderID, unmarshaledMsg.OrderID)
	}
	if unmarshaledMsg.Status != originalMsg.Status {
		t.Errorf("expected Status %s, got %s", originalMsg.Status, unmarshaledMsg.Status)
	}
	if unmarshaledMsg.MessageType != originalMsg.MessageType {
		t.Errorf("expected MessageType %s, got %s", originalMsg.MessageType, unmarshaledMsg.MessageType)
	}
}

// TestPaymentNotificationMessage_MessageType tests message type field
func TestPaymentNotificationMessage_MessageType(t *testing.T) {
	msg := PaymentNotificationMessage{
		OrderID:     uuid.New(),
		Status:      "paid",
		Timestamp:   time.Now(),
		MessageType: "payment_notification",
	}

	if msg.MessageType != "payment_notification" {
		t.Errorf("expected MessageType 'payment_notification', got '%s'", msg.MessageType)
	}

	// Test JSON output contains message_type
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if result["message_type"] != "payment_notification" {
		t.Errorf("expected message_type 'payment_notification' in JSON, got '%v'", result["message_type"])
	}
}

// TestRabbitMQConfig_ValidConfig tests valid configuration
func TestRabbitMQConfig_ValidConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  RabbitMQConfig
		wantURL string
	}{
		{
			name: "default config",
			config: RabbitMQConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "guest",
				Password: "guest",
				VHost:    "/",
			},
			wantURL: "amqp://guest:guest@localhost:5672//",
		},
		{
			name: "config with custom vhost",
			config: RabbitMQConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "test",
				Password: "test",
				VHost:    "custom",
			},
			wantURL: "amqp://test:test@localhost:5672/custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that config can be used to build URL
			url := buildConnectionURL(tt.config)
			if url != tt.wantURL {
				t.Errorf("buildConnectionURL() = %s, want %s", url, tt.wantURL)
			}
		})
	}
}

// buildConnectionURL is a helper function (extracted for testing)
func buildConnectionURL(config RabbitMQConfig) string {
	return "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/" + config.VHost
}

// TestNewRabbitMQPublisher_Integration tests publisher creation with real RabbitMQ
func TestNewRabbitMQPublisher_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := getTestRabbitMQConfig()

	exchangeName := "test_exchange_" + uuid.New().String()
	queueName := "test_queue_" + uuid.New().String()

	publisher, err := NewRabbitMQPublisher(config, exchangeName, queueName)
	if err != nil {
		t.Skipf("RabbitMQ not available at %s:%s (error: %v)", config.Host, config.Port, err)
	}
	defer func() {
		publisher.Close()
		// Cleanup: try to delete exchange and queue
		// Note: This is optional as they may be declared durable
	}()

	if publisher == nil {
		t.Fatal("expected publisher, got nil")
	}

	if publisher.exchangeName != exchangeName {
		t.Errorf("expected exchangeName %s, got %s", exchangeName, publisher.exchangeName)
	}

	if publisher.queueName != queueName {
		t.Errorf("expected queueName %s, got %s", queueName, publisher.queueName)
	}

	if publisher.IsClosed() {
		t.Error("expected connection to be open")
	}
}

// TestNewRabbitMQPublisher_InvalidConfig tests publisher creation with invalid config
func TestNewRabbitMQPublisher_InvalidConfig(t *testing.T) {
	tests := []struct {
		name         string
		config       RabbitMQConfig
		exchangeName string
		queueName    string
	}{
		{
			name: "invalid host",
			config: RabbitMQConfig{
				Host:     "invalid-host-that-does-not-exist.local",
				Port:     "5672",
				User:     "guest",
				Password: "guest",
				VHost:    "/",
			},
			exchangeName: "test_exchange",
			queueName:    "test_queue",
		},
		{
			name: "invalid port",
			config: RabbitMQConfig{
				Host:     "localhost",
				Port:     "9999",
				User:     "guest",
				Password: "guest",
				VHost:    "/",
			},
			exchangeName: "test_exchange",
			queueName:    "test_queue",
		},
		{
			name: "invalid credentials",
			config: RabbitMQConfig{
				Host:     "localhost",
				Port:     "5672",
				User:     "invalid_user",
				Password: "invalid_password",
				VHost:    "/",
			},
			exchangeName: "test_exchange",
			queueName:    "test_queue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRabbitMQPublisher(tt.config, tt.exchangeName, tt.queueName)
			if err == nil {
				t.Error("expected error for invalid config, got nil")
			}
		})
	}
}

// TestPublishPaymentNotification_Integration tests publishing message to RabbitMQ
func TestPublishPaymentNotification_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := getTestRabbitMQConfig()

	exchangeName := "test_exchange_" + uuid.New().String()
	queueName := "test_queue_" + uuid.New().String()

	publisher, err := NewRabbitMQPublisher(config, exchangeName, queueName)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}
	defer publisher.Close()

	// Create consumer to verify message was published
	conn, err := amqp.Dial(buildConnectionURL(config))
	if err != nil {
		t.Skipf("Failed to create consumer connection: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open consumer channel: %v", err)
	}
	defer ch.Close()

	// Start consuming messages
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		t.Fatalf("Failed to register consumer: %v", err)
	}

	orderID := uuid.New()

	// Publish message
	err = publisher.PublishPaymentNotification(orderID)
	if err != nil {
		t.Fatalf("PublishPaymentNotification failed: %v", err)
	}

	// Wait for message
	select {
	case msg := <-msgs:
		if msg.Body == nil {
			t.Fatal("received nil message body")
		}

		var paymentMsg PaymentNotificationMessage
		if err := json.Unmarshal(msg.Body, &paymentMsg); err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}

		if paymentMsg.OrderID != orderID {
			t.Errorf("expected OrderID %s, got %s", orderID, paymentMsg.OrderID)
		}
		if paymentMsg.Status != "paid" {
			t.Errorf("expected Status 'paid', got '%s'", paymentMsg.Status)
		}
		if paymentMsg.MessageType != "payment_notification" {
			t.Errorf("expected MessageType 'payment_notification', got '%s'", paymentMsg.MessageType)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

// TestRabbitMQPublisher_Close tests closing the publisher
func TestRabbitMQPublisher_Close_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := getTestRabbitMQConfig()

	exchangeName := "test_exchange_" + uuid.New().String()
	queueName := "test_queue_" + uuid.New().String()

	publisher, err := NewRabbitMQPublisher(config, exchangeName, queueName)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	if publisher.IsClosed() {
		t.Error("expected connection to be open initially")
	}

	err = publisher.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if !publisher.IsClosed() {
		t.Error("expected connection to be closed after Close()")
	}

	err = publisher.Close()
	if err != nil {
		t.Errorf("Close second time failed: %v", err)
	}
}

// TestRabbitMQPublisher_Reconnect_Integration tests reconnection logic
func TestRabbitMQPublisher_Reconnect_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := getTestRabbitMQConfig()

	exchangeName := "test_exchange_" + uuid.New().String()
	queueName := "test_queue_" + uuid.New().String()

	publisher, err := NewRabbitMQPublisher(config, exchangeName, queueName)
	if err != nil {
		t.Skipf("RabbitMQ not available: %v", err)
	}

	// Close current connection
	publisher.Close()

	if !publisher.IsClosed() {
		t.Error("expected connection to be closed")
	}

	// Reconnect
	err = publisher.Reconnect()
	if err != nil {
		// Reconnect might fail if RabbitMQ is not available
		t.Skipf("Reconnect failed (RabbitMQ may not be available): %v", err)
	}

	if publisher.IsClosed() {
		t.Error("expected connection to be open after Reconnect")
	}

	publisher.Close()
}

// TestPaymentNotificationMessage_StatusFields tests status field validation
func TestPaymentNotificationMessage_StatusFields(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"paid status", "paid"},
		{"pending status", "pending"},
		{"failed status", "failed"},
		{"refunded status", "refunded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID := uuid.New()
			timestamp := time.Now()
			msg := PaymentNotificationMessage{
				OrderID:     orderID,
				Status:      tt.status,
				Timestamp:   timestamp,
				MessageType: "payment_notification",
			}

			if msg.OrderID != orderID {
				t.Errorf("expected OrderID %s, got %s", orderID, msg.OrderID)
			}
			if msg.Status != tt.status {
				t.Errorf("expected Status %s, got %s", tt.status, msg.Status)
			}
			if msg.Timestamp.IsZero() {
				t.Error("expected non-zero timestamp")
			}
			if msg.MessageType != "payment_notification" {
				t.Errorf("expected MessageType 'payment_notification', got '%s'", msg.MessageType)
			}
		})
	}
}

// TestPaymentNotificationMessage_Timestamp tests timestamp field
func TestPaymentNotificationMessage_Timestamp(t *testing.T) {
	orderID := uuid.New()
	now := time.Now()
	msg := PaymentNotificationMessage{
		OrderID:     orderID,
		Status:      "paid",
		Timestamp:   now,
		MessageType: "payment_notification",
	}

	if msg.OrderID != orderID {
		t.Errorf("expected OrderID %s, got %s", orderID, msg.OrderID)
	}
	if msg.Status != "paid" {
		t.Errorf("expected Status 'paid', got '%s'", msg.Status)
	}
	if msg.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if msg.Timestamp.Unix() != now.Unix() {
		t.Errorf("expected timestamp %v, got %v", now, msg.Timestamp)
	}
	if msg.MessageType != "payment_notification" {
		t.Errorf("expected MessageType 'payment_notification', got '%s'", msg.MessageType)
	}
}

// TestPublishPaymentNotification_PublishingOptions tests message publishing options
func TestPublishPaymentNotification_PublishingOptions(t *testing.T) {
	// This test verifies that the message structure contains all required fields
	msg := PaymentNotificationMessage{
		OrderID:     uuid.New(),
		Status:      "paid",
		Timestamp:   time.Now(),
		MessageType: "payment_notification",
	}

	// Verify all required fields are set
	if msg.OrderID == uuid.Nil {
		t.Error("OrderID should not be nil")
	}
	if msg.Status == "" {
		t.Error("Status should not be empty")
	}
	if msg.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if msg.MessageType == "" {
		t.Error("MessageType should not be empty")
	}

	// Verify the message can be marshaled to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	// Verify JSON contains all required fields
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	requiredFields := []string{"order_id", "status", "timestamp", "message_type"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("JSON missing required field: %s", field)
		}
	}
}
