package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/domain"
	"github.com/nntch26/todo-app/backend/internal/middleware"
	"github.com/nntch26/todo-app/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// authRequest is the body for register and login
type authRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"` // bcrypt only uses the first 72 bytes
}

// userResponse is what clients see; never includes PasswordHash
type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email}
}

// Register
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req authRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}

	u, err := h.auth.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(toUserResponse(u))
}

// Login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req authRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}

	u, token, err := h.auth.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	setTokenCookie(c, token, time.Now().Add(service.TokenTTL))
	return c.JSON(toUserResponse(u))
}

// Logout expires the cookie
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	setTokenCookie(c, "", time.Unix(0, 0))
	return c.SendStatus(fiber.StatusNoContent)
}

// Me GET /api/auth/me (login required)
func (h *AuthHandler) Me(c fiber.Ctx) error {
	u, err := h.auth.Me(c.Context(), middleware.UserID(c))
	if err != nil {
		return err
	}
	return c.JSON(toUserResponse(u))
}

// set Token Cookie
func setTokenCookie(c fiber.Ctx, value string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     middleware.CookieName,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,                        // JavaScript can't read it
		Secure:   true,                        // HTTPS only; browsers treat http://localhost as secure too
		SameSite: fiber.CookieSameSiteLaxMode, // not sent on cross-site POSTs (basic CSRF protection)
	})
}
