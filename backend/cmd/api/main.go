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
	"github.com/nntch26/todo-app/backend/internal/middleware"
	"github.com/nntch26/todo-app/backend/internal/repository"
	"github.com/nntch26/todo-app/backend/internal/service"
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

	
	// Wire layers: repository (DB) -> service  -> handler (HTTP)
	// Auth
	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)
	requireAuth := middleware.RequireAuth(authService) // checks the login cookie

	// Todo
	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository)
	todoHandler := handler.NewTodoHandler(todoService)

	// ErrorHandler turns every returned error into a JSON response
	app := fiber.New(fiber.Config{ErrorHandler: handler.ErrorHandler})

	// Register routes
	api := app.Group("/api")
	api.Get("/health", handler.Health)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", requireAuth, authHandler.Logout) // middleware first: runs left to right
	auth.Get("/me", requireAuth, authHandler.Me)

	todos := api.Group("/todos", requireAuth) // every todo route needs login
	todos.Get("/", todoHandler.List)
	todos.Post("/", todoHandler.Create)
	todos.Patch("/:id", todoHandler.Update)
	todos.Delete("/:id", todoHandler.Delete)

	// Start server
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	// Wait for Ctrl+C / docker stop
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
