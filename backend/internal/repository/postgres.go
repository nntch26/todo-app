// Package repository is the data access layer (GORM + PostgreSQL).
package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// Connect opens the DB
func Connect(dsn string) (*gorm.DB, error) {
	// TranslateError: turn Postgres errors into gorm.ErrDuplicatedKey etc.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         newSafeLogger(),
	})
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

// newSafeLogger logs SQL as "WHERE email = $1" (no values), so emails/hashes never reach the logs
func newSafeLogger() logger.Interface {
	return logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true, // "not found" is normal (e.g. login with unknown email)
		ParameterizedQueries:      true,
	})
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
