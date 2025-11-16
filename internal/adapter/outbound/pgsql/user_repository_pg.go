package pgsql

import (
	"context"
	"errors"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryPG implements UserRepository interface for PostgreSQL using GORM
type UserRepositoryPG struct {
	db *gorm.DB
}

// NewUserRepositoryPG creates a new PostgreSQL user repository
func NewUserRepositoryPG(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryPG{db: db}
}

// Create creates a new user
func (r *UserRepositoryPG) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindByID finds a user by ID
func (r *UserRepositoryPG) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryPG) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepositoryPG) Update(ctx context.Context, user *domain.User) error {
	result := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":       user.Name,
			"updated_at": user.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user (soft delete using GORM)
func (r *UserRepositoryPG) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepositoryPG) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Count counts total users
func (r *UserRepositoryPG) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepositoryPG) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
