package event

import (
	"context"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/port/outbound/broker"
)

// UserEventPublisher publishes user domain events
type UserEventPublisher struct {
	broker broker.MessageBroker
}

// NewUserEventPublisher creates a new user event publisher
func NewUserEventPublisher(broker broker.MessageBroker) *UserEventPublisher {
	return &UserEventPublisher{
		broker: broker,
	}
}

// PublishUserCreated publishes a user created event
func (p *UserEventPublisher) PublishUserCreated(ctx context.Context, event *domain.UserCreatedEvent) error {
	if p.broker == nil {
		return nil // Gracefully handle when broker is disabled
	}

	if err := p.broker.Publish(ctx, "user.created", event); err != nil {
		return fmt.Errorf("failed to publish user created event: %w", err)
	}

	return nil
}

// PublishUserUpdated publishes a user updated event
func (p *UserEventPublisher) PublishUserUpdated(ctx context.Context, event *domain.UserUpdatedEvent) error {
	if p.broker == nil {
		return nil
	}

	if err := p.broker.Publish(ctx, "user.updated", event); err != nil {
		return fmt.Errorf("failed to publish user updated event: %w", err)
	}

	return nil
}

// PublishUserDeleted publishes a user deleted event
func (p *UserEventPublisher) PublishUserDeleted(ctx context.Context, event *domain.UserDeletedEvent) error {
	if p.broker == nil {
		return nil
	}

	if err := p.broker.Publish(ctx, "user.deleted", event); err != nil {
		return fmt.Errorf("failed to publish user deleted event: %w", err)
	}

	return nil
}

// PublishUserLoggedIn publishes a user logged in event
func (p *UserEventPublisher) PublishUserLoggedIn(ctx context.Context, event *domain.UserLoggedInEvent) error {
	if p.broker == nil {
		return nil
	}

	if err := p.broker.Publish(ctx, "user.logged_in", event); err != nil {
		return fmt.Errorf("failed to publish user logged in event: %w", err)
	}

	return nil
}
