package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PaymentNotificationMessage represents the message sent to RabbitMQ
type PaymentNotificationMessage struct {
	OrderID     uuid.UUID `json:"order_id"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"message_type"` // "payment_notification"
}

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	VHost    string // Virtual host, default "/"
}

// Connection interface defines methods we use from amqp.Connection
type Connection interface {
	Close() error
	IsClosed() bool
	Channel() (Channel, error)
}

// Channel interface defines methods we use from amqp.Channel
type Channel interface {
	Close() error
	ExchangeDeclare(exchange, key string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

// amqpConnection wraps amqp.Connection to implement Connection interface
type amqpConnection struct{ *amqp.Connection }

func (a amqpConnection) Channel() (Channel, error) {
	ch, err := a.Connection.Channel()
	if err != nil {
		return nil, err
	}
	return amqpChannel{ch}, nil
}

// amqpChannel wraps amqp.Channel to implement Channel interface
type amqpChannel struct{ *amqp.Channel }

// RabbitMQPublisher handles publishing messages to RabbitMQ
type RabbitMQPublisher struct {
	conn         Connection
	channel      Channel
	exchangeName string
	queueName    string
	config       RabbitMQConfig
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher
func NewRabbitMQPublisher(config RabbitMQConfig, exchangeName, queueName string) (*RabbitMQPublisher, error) {
	publisher := &RabbitMQPublisher{
		config:       config,
		exchangeName: exchangeName,
		queueName:    queueName,
	}

	if err := publisher.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return publisher, nil
}

// connect establishes connection to RabbitMQ and sets up exchange/queue
func (p *RabbitMQPublisher) connect() error {
	// Construct connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		p.config.User,
		p.config.Password,
		p.config.Host,
		p.config.Port,
		p.config.VHost,
	)

	var err error
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}
	p.conn = amqpConnection{conn}

	// Create channel
	rawChannel, err := p.conn.Channel()
	if err != nil {
		p.conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}
	p.channel = rawChannel

	// Declare exchange (fanout for push notifications to multiple consumers)
	if err := p.channel.ExchangeDeclare(
		p.exchangeName, // name
		"fanout",       // type (fanout, direct, topic, headers)
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	); err != nil {
		p.channel.Close()
		p.conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	if _, err := p.channel.QueueDeclare(
		p.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		p.channel.Close()
		p.conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	if err := p.channel.QueueBind(
		p.queueName,    // queue name
		"",             // routing key (empty for fanout)
		p.exchangeName, // exchange
		false,          // no-wait
		nil,            // arguments
	); err != nil {
		p.channel.Close()
		p.conn.Close()
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Printf("Connected to RabbitMQ: exchange=%s, queue=%s", p.exchangeName, p.queueName)
	return nil
}

// PublishPaymentNotification publishes a payment notification message
func (p *RabbitMQPublisher) PublishPaymentNotification(orderID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := PaymentNotificationMessage{
		OrderID:     orderID,
		Status:      "paid",
		Timestamp:   time.Now(),
		MessageType: "payment_notification",
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish to exchange (will be routed to bound queues)
	err = p.channel.PublishWithContext(
		ctx,
		p.exchangeName, // exchange
		"",             // routing key (ignored for fanout)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published payment notification for order %s", orderID)
	return nil
}

// Reconnect attempts to reconnect to RabbitMQ
func (p *RabbitMQPublisher) Reconnect() error {
	if p.conn != nil && !p.conn.IsClosed() {
		p.conn.Close()
	}
	return p.connect()
}

// Close closes the RabbitMQ connection (idempotent - safe to call multiple times)
func (p *RabbitMQPublisher) Close() error {
	var err error
	if p.channel != nil && !p.conn.IsClosed() {
		if closeErr := p.channel.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if p.conn != nil && !p.conn.IsClosed() {
		if closeErr := p.conn.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

// IsClosed returns true if the connection is closed
func (p *RabbitMQPublisher) IsClosed() bool {
	return p.conn == nil || p.conn.IsClosed()
}
