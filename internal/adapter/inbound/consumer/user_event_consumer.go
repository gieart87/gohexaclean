package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/port/outbound/broker"
)

// UserEventConsumer consumes user domain events
type UserEventConsumer struct {
	broker broker.MessageBroker
}

// NewUserEventConsumer creates a new user event consumer
func NewUserEventConsumer(broker broker.MessageBroker) *UserEventConsumer {
	return &UserEventConsumer{
		broker: broker,
	}
}

// Start starts consuming user events
func (c *UserEventConsumer) Start(ctx context.Context) error {
	if c.broker == nil {
		return nil // Gracefully handle when broker is disabled
	}

	// Subscribe to user created events
	if err := c.broker.Subscribe(ctx, "user.created", c.handleUserCreated); err != nil {
		return fmt.Errorf("failed to subscribe to user.created: %w", err)
	}

	// Subscribe to user updated events
	if err := c.broker.Subscribe(ctx, "user.updated", c.handleUserUpdated); err != nil {
		return fmt.Errorf("failed to subscribe to user.updated: %w", err)
	}

	// Subscribe to user deleted events
	if err := c.broker.Subscribe(ctx, "user.deleted", c.handleUserDeleted); err != nil {
		return fmt.Errorf("failed to subscribe to user.deleted: %w", err)
	}

	// Subscribe to user logged in events
	if err := c.broker.Subscribe(ctx, "user.logged_in", c.handleUserLoggedIn); err != nil {
		return fmt.Errorf("failed to subscribe to user.logged_in: %w", err)
	}

	return nil
}

// Stop stops consuming user events
func (c *UserEventConsumer) Stop() error {
	if c.broker == nil {
		return nil
	}

	topics := []string{"user.created", "user.updated", "user.deleted", "user.logged_in"}
	for _, topic := range topics {
		if err := c.broker.Unsubscribe(topic); err != nil {
			log.Printf("failed to unsubscribe from %s: %v", topic, err)
		}
	}

	return nil
}

// handleUserCreated handles user created events
func (c *UserEventConsumer) handleUserCreated(ctx context.Context, message []byte) error {
	var event domain.UserCreatedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user created event: %w", err)
	}

	log.Printf("[EVENT] User Created: ID=%s, Email=%s, Name=%s, At=%s",
		event.AggregateID(), event.Email, event.Name, event.OccurredAt())

	// Add your business logic here
	// For example:
	// - Send welcome email
	// - Create user profile in another service
	// - Update analytics
	// - Send notification

	return nil
}

// handleUserUpdated handles user updated events
func (c *UserEventConsumer) handleUserUpdated(ctx context.Context, message []byte) error {
	var event domain.UserUpdatedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user updated event: %w", err)
	}

	log.Printf("[EVENT] User Updated: ID=%s, Name=%s, At=%s",
		event.AggregateID(), event.Name, event.OccurredAt())

	// Add your business logic here
	// For example:
	// - Sync with external systems
	// - Invalidate cache
	// - Update search index

	return nil
}

// handleUserDeleted handles user deleted events
func (c *UserEventConsumer) handleUserDeleted(ctx context.Context, message []byte) error {
	var event domain.UserDeletedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user deleted event: %w", err)
	}

	log.Printf("[EVENT] User Deleted: ID=%s, At=%s",
		event.AggregateID(), event.OccurredAt())

	// Add your business logic here
	// For example:
	// - Clean up user data
	// - Archive user information
	// - Send deletion notification
	// - Remove from external systems

	return nil
}

// handleUserLoggedIn handles user logged in events
func (c *UserEventConsumer) handleUserLoggedIn(ctx context.Context, message []byte) error {
	var event domain.UserLoggedInEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user logged in event: %w", err)
	}

	log.Printf("[EVENT] User Logged In: ID=%s, Email=%s, At=%s",
		event.AggregateID(), event.Email, event.OccurredAt())

	// Add your business logic here
	// For example:
	// - Track login analytics
	// - Update last login timestamp
	// - Send login notification
	// - Check for suspicious activity

	return nil
}
