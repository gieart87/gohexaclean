package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user domain model (entity)
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email     string         `gorm:"uniqueIndex;not null;size:255"`
	Name      string         `gorm:"not null;size:255"`
	Password  string         `gorm:"not null;size:255"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user entity
func NewUser(email, name, password string) *User {
	return &User{
		ID:       uuid.New(),
		Email:    email,
		Name:     name,
		Password: password,
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

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name string) {
	u.Name = name
}
