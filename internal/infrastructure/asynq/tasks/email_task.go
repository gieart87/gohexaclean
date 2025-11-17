package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailWelcome = "email:welcome"
)

// EmailWelcomePayload represents the payload for welcome email task
type EmailWelcomePayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

// NewEmailWelcomeTask creates a new task to send welcome email
func NewEmailWelcomeTask(userID, email, name string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailWelcomePayload{
		UserID: userID,
		Email:  email,
		Name:   name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmailWelcome, payload), nil
}

// HandleEmailWelcomeTask processes the welcome email task
func HandleEmailWelcomeTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailWelcomePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// TODO: Implement actual email sending logic here
	// For now, we'll just log it
	log.Printf("Sending welcome email to %s (%s) for user %s", payload.Name, payload.Email, payload.UserID)

	// Simulate email sending
	// In production, you would use an email service like SendGrid, AWS SES, etc.
	log.Printf("Welcome email sent successfully to %s", payload.Email)

	return nil
}
