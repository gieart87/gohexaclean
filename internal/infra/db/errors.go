package db

import "errors"

// Database infrastructure errors
var (
	ErrDBConnection     = errors.New("database connection failed")
	ErrDBTimeout        = errors.New("database operation timeout")
	ErrDBTransaction    = errors.New("database transaction failed")
	ErrDBMigration      = errors.New("database migration failed")
	ErrDBRecordNotFound = errors.New("record not found in database")
	ErrDBDuplicateKey   = errors.New("duplicate key violation")
	ErrDBConstraint     = errors.New("database constraint violation")
)
