package broker

import (
	"fmt"

	"github.com/gieart87/gohexaclean/internal/adapter/outbound/rabbitmq"
	"github.com/gieart87/gohexaclean/internal/infra/config"
	"github.com/gieart87/gohexaclean/internal/port/outbound/broker"
)

// NewMessageBroker creates a new message broker based on configuration
func NewMessageBroker(cfg *config.BrokerConfig) (broker.MessageBroker, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("message broker is disabled")
	}

	switch cfg.Type {
	case "rabbitmq":
		return rabbitmq.NewRabbitMQBroker(&cfg.RabbitMQ), nil
	// Future broker implementations can be added here
	// case "kafka":
	//     return kafka.NewKafkaBroker(&cfg.Kafka), nil
	// case "pubsub":
	//     return pubsub.NewPubSubBroker(&cfg.PubSub), nil
	// case "nats":
	//     return nats.NewNATSBroker(&cfg.NATS), nil
	default:
		return nil, fmt.Errorf("unsupported broker type: %s", cfg.Type)
	}
}
