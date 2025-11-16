package pgsql

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gieart87/gohexaclean/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	return db, mock
}

func TestUserRepositoryPG_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	user := &domain.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(user.Email, user.Name, user.Password, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_FindByID(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at", "updated_at", "deleted_at"}).
		AddRow(userID, "test@example.com", "Test User", "hashedpassword", now, now, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT`)).
		WithArgs(userID, 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_FindByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	userID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT`)).
		WithArgs(userID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByID(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_FindByEmail(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	userID := uuid.New()
	email := "test@example.com"
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at", "updated_at", "deleted_at"}).
		AddRow(userID, email, "Test User", "hashedpassword", now, now, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_FindByEmail_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	email := "notfound@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT`)).
		WithArgs(email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	user := &domain.User{
		ID:        uuid.New(),
		Name:      "Updated Name",
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs(user.Name, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_Update_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	user := &domain.User{
		ID:        uuid.New(),
		Name:      "Updated Name",
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs(user.Name, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_Delete_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	userID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE id = $2 AND "users"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(context.Background(), userID)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_List(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "email", "name", "password", "created_at", "updated_at", "deleted_at"}).
		AddRow(uuid.New(), "user1@example.com", "User 1", "pass1", now, now, nil).
		AddRow(uuid.New(), "user2@example.com", "User 2", "pass2", now, now, nil)

	// GORM doesn't add OFFSET when it's 0
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(rows)

	users, err := repo.List(context.Background(), 0, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_Count(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	count, err := repo.Count(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_ExistsByEmail(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryPG_ExistsByEmail_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepositoryPG(db)

	email := "notfound@example.com"
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL`)).
		WithArgs(email).
		WillReturnRows(rows)

	exists, err := repo.ExistsByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}
