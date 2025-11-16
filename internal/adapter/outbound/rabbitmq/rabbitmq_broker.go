package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/port/outbound/broker"
)

// RabbitMQBroker implements the MessageBroker interface for RabbitMQ
type RabbitMQBroker struct {
	config     *config.RabbitMQConfig
	conn       *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	connected  bool
	reconnecting bool
	subscriptions map[string]*subscription
	done       chan struct{}
}

type subscription struct {
	queue   string
	handler broker.MessageHandler
	cancel  context.CancelFunc
}

// NewRabbitMQBroker creates a new RabbitMQ message broker
func NewRabbitMQBroker(cfg *config.RabbitMQConfig) *RabbitMQBroker {
	return &RabbitMQBroker{
		config:        cfg,
		subscriptions: make(map[string]*subscription),
		done:          make(chan struct{}),
	}
}

// Connect establishes connection to RabbitMQ
func (r *RabbitMQBroker) Connect(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		return nil
	}

	connConfig := amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	}

	if r.config.ConnectionName != "" {
		connConfig.Properties = amqp.Table{
			"connection_name": r.config.ConnectionName,
		}
	}

	conn, err := amqp.DialConfig(r.config.GetAMQPURL(), connConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS (prefetch count)
	if err := ch.Qos(r.config.PrefetchCount, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declare exchange
	if r.config.Exchange != "" {
		exchangeType := r.config.ExchangeType
		if exchangeType == "" {
			exchangeType = "topic"
		}

		err = ch.ExchangeDeclare(
			r.config.Exchange,
			exchangeType,
			true,  // durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to declare exchange: %w", err)
		}
	}

	r.conn = conn
	r.channel = ch
	r.connected = true

	// Monitor connection
	go r.monitorConnection()

	return nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQBroker) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.connected {
		return nil
	}

	close(r.done)
	r.connected = false

	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}

	return nil
}

// Health checks the health of RabbitMQ connection
func (r *RabbitMQBroker) Health() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.connected || r.conn == nil || r.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is not healthy")
	}

	return nil
}

// Publish publishes an event to RabbitMQ
func (r *RabbitMQBroker) Publish(ctx context.Context, topic string, event domain.Event) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	exchange := r.config.Exchange
	if exchange == "" {
		exchange = "amq.topic"
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Transient,
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    event.OccurredAt(),
		MessageId:    event.EventID(),
		Type:         event.EventType(),
	}

	if r.config.Persistent {
		msg.DeliveryMode = amqp.Persistent
	}

	err = r.channel.PublishWithContext(
		ctx,
		exchange,
		topic,
		false, // mandatory
		false, // immediate
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishBatch publishes multiple events in a batch
func (r *RabbitMQBroker) PublishBatch(ctx context.Context, topic string, events []domain.Event) error {
	for _, event := range events {
		if err := r.Publish(ctx, topic, event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe subscribes to a topic and handles incoming messages
func (r *RabbitMQBroker) Subscribe(ctx context.Context, topic string, handler broker.MessageHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.connected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Check if already subscribed
	if _, exists := r.subscriptions[topic]; exists {
		return fmt.Errorf("already subscribed to topic: %s", topic)
	}

	// Create queue name
	queueName := r.config.QueuePrefix + topic
	if r.config.QueuePrefix == "" {
		queueName = topic
	}

	// Declare queue
	queue, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	exchange := r.config.Exchange
	if exchange == "" {
		exchange = "amq.topic"
	}

	err = r.channel.QueueBind(
		queue.Name,
		topic,
		exchange,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Start consuming
	msgs, err := r.channel.Consume(
		queue.Name,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// Create subscription context
	subCtx, cancel := context.WithCancel(ctx)

	// Store subscription
	r.subscriptions[topic] = &subscription{
		queue:   queue.Name,
		handler: handler,
		cancel:  cancel,
	}

	// Process messages
	go r.processMessages(subCtx, topic, msgs)

	return nil
}

// Unsubscribe unsubscribes from a topic
func (r *RabbitMQBroker) Unsubscribe(topic string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, exists := r.subscriptions[topic]
	if !exists {
		return fmt.Errorf("not subscribed to topic: %s", topic)
	}

	// Cancel subscription
	sub.cancel()

	// Remove from subscriptions
	delete(r.subscriptions, topic)

	return nil
}

// processMessages processes incoming messages from a queue
func (r *RabbitMQBroker) processMessages(ctx context.Context, topic string, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			r.mu.RLock()
			sub, exists := r.subscriptions[topic]
			r.mu.RUnlock()

			if !exists {
				msg.Nack(false, true) // Requeue if subscription was removed
				continue
			}

			// Handle message
			if err := sub.handler(ctx, msg.Body); err != nil {
				// Nack and requeue on error
				msg.Nack(false, true)
			} else {
				// Ack on success
				msg.Ack(false)
			}
		}
	}
}

// monitorConnection monitors the connection and attempts to reconnect
func (r *RabbitMQBroker) monitorConnection() {
	closeChan := make(chan *amqp.Error)
	r.conn.NotifyClose(closeChan)

	select {
	case <-r.done:
		return
	case err := <-closeChan:
		if err != nil {
			r.handleDisconnect()
		}
	}
}

// handleDisconnect handles connection loss and attempts reconnection
func (r *RabbitMQBroker) handleDisconnect() {
	r.mu.Lock()
	if r.reconnecting {
		r.mu.Unlock()
		return
	}
	r.reconnecting = true
	r.connected = false
	r.mu.Unlock()

	attempts := 0
	maxAttempts := r.config.MaxReconnect
	if maxAttempts == 0 {
		maxAttempts = 10
	}

	delay := r.config.ReconnectDelay
	if delay == 0 {
		delay = 5 * time.Second
	}

	for attempts < maxAttempts {
		select {
		case <-r.done:
			return
		case <-time.After(delay):
			if err := r.Connect(context.Background()); err == nil {
				r.mu.Lock()
				r.reconnecting = false
				r.mu.Unlock()

				// Resubscribe to all topics
				r.resubscribeAll()
				return
			}
			attempts++
		}
	}

	r.mu.Lock()
	r.reconnecting = false
	r.mu.Unlock()
}

// resubscribeAll resubscribes to all topics after reconnection
func (r *RabbitMQBroker) resubscribeAll() {
	r.mu.RLock()
	topics := make([]string, 0, len(r.subscriptions))
	handlers := make(map[string]broker.MessageHandler)

	for topic, sub := range r.subscriptions {
		topics = append(topics, topic)
		handlers[topic] = sub.handler
	}
	r.mu.RUnlock()

	// Clear old subscriptions
	r.mu.Lock()
	r.subscriptions = make(map[string]*subscription)
	r.mu.Unlock()

	// Resubscribe
	for _, topic := range topics {
		r.Subscribe(context.Background(), topic, handlers[topic])
	}
}
