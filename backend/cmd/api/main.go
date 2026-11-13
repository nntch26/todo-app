package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/config"
	"github.com/nntch26/todo-app/backend/internal/handler"
	"github.com/nntch26/todo-app/backend/internal/repository"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	// Connect database + migrate tables
	db, err := repository.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("database: %v", err)
	}

	app := fiber.New()

	// Register routes
	app.Get("/api/health", handler.Health)

	// Start server
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	//
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}

	// Close database
	if err := repository.Close(db); err != nil {
		log.Fatalf("close database: %v", err)
	}
}
