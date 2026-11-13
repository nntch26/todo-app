// Package repository is the data access layer (GORM + PostgreSQL).
package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// Connect opens the DB
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Migrate creates/updates tables.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}, &domain.Todo{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}

// Close releases the connection pool
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	return sqlDB.Close()
}
