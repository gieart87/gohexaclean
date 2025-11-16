package domain

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event
type Event interface {
	EventType() string
	EventID() string
	OccurredAt() time.Time
	AggregateID() string
}

// BaseEvent provides common event fields
type BaseEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
	AggregateId string    `json:"aggregate_id"`
}

func (e BaseEvent) EventType() string {
	return e.Type
}

func (e BaseEvent) EventID() string {
	return e.ID
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (e BaseEvent) AggregateID() string {
	return e.AggregateId
}

// User Domain Events

// UserCreatedEvent is published when a new user is created
type UserCreatedEvent struct {
	BaseEvent
	Email string `json:"email"`
	Name  string `json:"name"`
}

func NewUserCreatedEvent(userID uuid.UUID, email, name string) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.created",
			Timestamp:   time.Now(),
			AggregateId: userID.String(),
		},
		Email: email,
		Name:  name,
	}
}

// UserUpdatedEvent is published when a user is updated
type UserUpdatedEvent struct {
	BaseEvent
	Name string `json:"name"`
}

func NewUserUpdatedEvent(userID uuid.UUID, name string) *UserUpdatedEvent {
	return &UserUpdatedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.updated",
			Timestamp:   time.Now(),
			AggregateId: userID.String(),
		},
		Name: name,
	}
}

// UserDeletedEvent is published when a user is deleted
type UserDeletedEvent struct {
	BaseEvent
}

func NewUserDeletedEvent(userID uuid.UUID) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.deleted",
			Timestamp:   time.Now(),
			AggregateId: userID.String(),
		},
	}
}

// UserLoggedInEvent is published when a user logs in
type UserLoggedInEvent struct {
	BaseEvent
	Email string `json:"email"`
}

func NewUserLoggedInEvent(userID uuid.UUID, email string) *UserLoggedInEvent {
	return &UserLoggedInEvent{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			Type:        "user.logged_in",
			Timestamp:   time.Now(),
			AggregateId: userID.String(),
		},
		Email: email,
	}
}
