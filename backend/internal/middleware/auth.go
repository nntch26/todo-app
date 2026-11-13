// Package middleware holds Fiber middlewares.
package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/domain"
	"github.com/nntch26/todo-app/backend/internal/service"
)

// CookieName is the httpOnly cookie that carries the JWT
const CookieName = "token"

// userIDKey is unexported so no other package can overwrite it in Locals
type userIDKey struct{}

// RequireAuth rejects requests without a valid token cookie, and stores the user ID for handlers
func RequireAuth(auth *service.AuthService) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, err := auth.ParseToken(c.Cookies(CookieName))
		if err != nil {
			return domain.ErrInvalidToken
		}

		c.Locals(userIDKey{}, userID)
		return c.Next()
	}
}

// UserID returns the logged-in user's ID (only valid on routes behind RequireAuth)
func UserID(c fiber.Ctx) string {
	id, _ := c.Locals(userIDKey{}).(string)
	return id
}
