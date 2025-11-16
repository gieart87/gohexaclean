# Message Broker Integration

This application supports flexible message broker integration with a pluggable architecture that allows easy switching between different broker implementations.

## Supported Brokers

Currently supported:
- **RabbitMQ** - Full implementation with auto-reconnect

Planned for future:
- **Apache Kafka**
- **Google Cloud Pub/Sub**
- **NATS**
- **AWS SQS/SNS**

## Architecture

The message broker implementation follows hexagonal architecture principles:

```
├── internal/
│   ├── domain/
│   │   └── events.go                    # Domain events (UserCreatedEvent, etc.)
│   ├── port/outbound/broker/
│   │   └── message_broker.go            # Port interfaces
│   ├── adapter/outbound/
│   │   ├── rabbitmq/
│   │   │   └── rabbitmq_broker.go       # RabbitMQ adapter
│   │   └── event/
│   │       └── user_event_publisher.go   # Event publisher
│   ├── adapter/inbound/consumer/
│   │   └── user_event_consumer.go       # Event consumer
│   └── infra/broker/
│       └── factory.go                   # Broker factory
```

## Configuration

### Environment Variables

```bash
# Enable/disable message broker
BROKER_ENABLED=true
BROKER_TYPE=rabbitmq

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
# Or individually:
RABBITMQ_HOST=localhost
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
```

### YAML Configuration

Add this to your `config.yaml`:

```yaml
broker:
  enabled: true
  type: rabbitmq  # rabbitmq, kafka, pubsub, nats
  rabbitmq:
    host: localhost
    port: 5672
    user: guest
    password: guest
    vhost: /
    exchange: user_events
    exchange_type: topic
    queue_prefix: gohexaclean_
    prefetch_count: 10
    reconnect_delay: 5s
    max_reconnect: 10
    persistent: true
    connection_name: gohexaclean-service
```

## Domain Events

The application publishes the following domain events:

### User Events

| Event Type | Routing Key | Description |
|------------|-------------|-------------|
| `user.created` | `user.created` | Published when a new user is created |
| `user.updated` | `user.updated` | Published when a user profile is updated |
| `user.deleted` | `user.deleted` | Published when a user is soft-deleted |
| `user.logged_in` | `user.logged_in` | Published when a user successfully logs in |

### Event Structure

All events implement the `domain.Event` interface:

```go
type Event interface {
    EventType() string
    EventID() string
    OccurredAt() time.Time
    AggregateID() string
}
```

Example `UserCreatedEvent`:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "user.created",
  "timestamp": "2025-11-16T10:30:00Z",
  "aggregate_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "name": "John Doe"
}
```

## Usage

### Publishing Events

Events are automatically published by the `UserService` when domain actions occur:

```go
// In UserService.CreateUser
event := domain.NewUserCreatedEvent(user.ID, user.Email, user.Name)
s.eventPublisher.PublishUserCreated(ctx, event)
```

### Consuming Events

The `UserEventConsumer` automatically subscribes to all user events. You can customize the handlers:

```go
// In user_event_consumer.go
func (c *UserEventConsumer) handleUserCreated(ctx context.Context, message []byte) error {
    var event domain.UserCreatedEvent
    if err := json.Unmarshal(message, &event); err != nil {
        return err
    }

    // Your business logic here
    // - Send welcome email
    // - Create user profile
    // - Update analytics

    return nil
}
```

### Adding Custom Event Handlers

1. Subscribe to events in your consumer:

```go
func (c *CustomConsumer) Start(ctx context.Context) error {
    return c.broker.Subscribe(ctx, "user.created", c.handleUserCreated)
}
```

2. Implement the handler:

```go
func (c *CustomConsumer) handleUserCreated(ctx context.Context, message []byte) error {
    // Process the message
    return nil
}
```

## Running with RabbitMQ

### Using Docker Compose

Create `docker-compose.rabbitmq.yml`:

```yaml
version: '3.8'

services:
  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: gohexaclean-rabbitmq
    ports:
      - "5672:5672"      # AMQP port
      - "15672:15672"    # Management UI
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  rabbitmq_data:
```

Start RabbitMQ:

```bash
docker-compose -f docker-compose.rabbitmq.yml up -d
```

Access RabbitMQ Management UI:
- URL: http://localhost:15672
- Username: guest
- Password: guest

### Testing Events

1. Start the application:

```bash
go run cmd/http/main.go
```

2. Register a new user:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'
```

3. Check RabbitMQ Management UI to see the published event

4. Check application logs to see the consumed event:

```
[EVENT] User Created: ID=xxx, Email=test@example.com, Name=Test User
```

## Error Handling

The message broker implementation includes:

- **Graceful degradation**: If the broker is disabled or fails, the application continues to work
- **Auto-reconnection**: Automatic reconnection with exponential backoff
- **Message acknowledgment**: Messages are acknowledged only after successful processing
- **Requeue on error**: Failed messages are requeued for retry

## Adding New Broker Implementations

To add support for a new broker (e.g., Kafka):

1. Create adapter implementation:

```go
// internal/adapter/outbound/kafka/kafka_broker.go
type KafkaBroker struct {
    // ...
}

func (k *KafkaBroker) Publish(ctx context.Context, topic string, event domain.Event) error {
    // Implement Kafka publishing
}

func (k *KafkaBroker) Subscribe(ctx context.Context, topic string, handler broker.MessageHandler) error {
    // Implement Kafka subscription
}
```

2. Add configuration:

```go
// internal/infra/config/config.go
type KafkaConfig struct {
    Brokers []string `yaml:"brokers"`
    // ...
}
```

3. Register in factory:

```go
// internal/infra/broker/factory.go
case "kafka":
    return kafka.NewKafkaBroker(&cfg.Kafka), nil
```

## Best Practices

1. **Event Design**:
   - Keep events immutable
   - Include all necessary data in the event
   - Use clear, descriptive event types

2. **Error Handling**:
   - Always handle errors in consumers
   - Log errors for debugging
   - Use dead letter queues for failed messages

3. **Performance**:
   - Use batch publishing when possible
   - Adjust prefetch count based on workload
   - Monitor queue depths

4. **Testing**:
   - Test with broker disabled (graceful degradation)
   - Test reconnection scenarios
   - Use mocks for unit tests

## Monitoring

Monitor these metrics:

- Message publish rate
- Message consumption rate
- Queue depth
- Consumer lag
- Connection status
- Error rate

## Troubleshooting

### Events not being published

1. Check if broker is enabled: `BROKER_ENABLED=true`
2. Verify connection settings
3. Check application logs for connection errors

### Events not being consumed

1. Verify consumer is started
2. Check queue bindings in RabbitMQ Management UI
3. Verify exchange and routing key configuration

### Connection issues

1. Check RabbitMQ is running: `docker ps`
2. Verify network connectivity
3. Check credentials
4. Review reconnection logs
