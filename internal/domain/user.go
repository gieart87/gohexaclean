package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user domain model (entity)
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Password  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user entity
func NewUser(email, name, password string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Password:  password,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate performs domain-level validation
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.Name == "" {
		return ErrNameRequired
	}
	if u.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// Deactivate marks user as inactive
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate marks user as active
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}
