package repository

import (
	"context"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines the outbound port for user repository
// This interface will be implemented by the database adapter
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*domain.User, error)
	Count(ctx context.Context) (int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
