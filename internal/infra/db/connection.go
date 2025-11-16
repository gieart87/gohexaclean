package db

import (
	"fmt"
	"time"

	"github.com/gieart87/gohexaclean/internal/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormConnection creates a new GORM database connection
func NewGormConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// Configure GORM logger (silent by default)
	gormLogger := logger.Default.LogMode(logger.Silent)

	// Open GORM connection
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Close closes the GORM database connection
func Close(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
