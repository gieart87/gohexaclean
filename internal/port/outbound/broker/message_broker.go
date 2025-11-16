package broker

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/domain"
)

// MessageBroker is the main interface for message broker operations
type MessageBroker interface {
	Publisher
	Consumer
	Connect(ctx context.Context) error
	Close() error
	Health() error
}

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(ctx context.Context, topic string, event domain.Event) error
	PublishBatch(ctx context.Context, topic string, events []domain.Event) error
}

// Consumer defines the interface for consuming messages
type Consumer interface {
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Unsubscribe(topic string) error
}

// MessageHandler is a function type for handling consumed messages
type MessageHandler func(ctx context.Context, message []byte) error

// Message represents a message consumed from the broker
type Message struct {
	ID        string
	Topic     string
	Body      []byte
	Timestamp int64
	Headers   map[string]string
}

// PublishOptions provides options for publishing messages
type PublishOptions struct {
	Topic       string
	Priority    uint8
	ContentType string
	Headers     map[string]string
	Persistent  bool
}
