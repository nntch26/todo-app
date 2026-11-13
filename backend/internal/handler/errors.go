package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// ErrorHandler turns every error into {"error": "..."}; internal details are logged, never sent to users
func ErrorHandler(c fiber.Ctx, err error) error {
	status, message := fiber.StatusInternalServerError, "internal server error"

	var fe *fiber.Error
	switch {
	case errors.As(err, &fe): // errors we created with fiber.NewError (400, 401, 404 route...)
		status, message = fe.Code, fe.Message
	case errors.Is(err, domain.ErrNotFound):
		status, message = fiber.StatusNotFound, "not found"
	case errors.Is(err, domain.ErrEmailTaken):
		status, message = fiber.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrInvalidToken):
		status, message = fiber.StatusUnauthorized, err.Error()
	default:
		log.Printf("internal error: %s %s: %v", c.Method(), c.Path(), err)
	}

	return c.Status(status).JSON(fiber.Map{"error": message})
}
