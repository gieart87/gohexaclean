package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/gieart87/gohexaclean/internal/port/outbound/repository"
	"github.com/google/uuid"
)

// UserRepositoryPG implements UserRepository interface for PostgreSQL
type UserRepositoryPG struct {
	db *sql.DB
}

// NewUserRepositoryPG creates a new PostgreSQL user repository
func NewUserRepositoryPG(db *sql.DB) repository.UserRepository {
	return &UserRepositoryPG{db: db}
}

// Create creates a new user
func (r *UserRepositoryPG) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, name, password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.Password,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID finds a user by ID
func (r *UserRepositoryPG) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepositoryPG) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepositoryPG) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, is_active = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *UserRepositoryPG) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepositoryPG) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	query := `
		SELECT id, email, name, password, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Count counts total users
func (r *UserRepositoryPG) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepositoryPG) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
