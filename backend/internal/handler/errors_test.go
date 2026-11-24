package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// doRequest sends a request to app and returns status + body
func doRequest(t *testing.T, app *fiber.App, method, path, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, string(b)
}

// TestErrorHandler checks status mapping and that internal details never reach the client
func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{"internal error is hidden", errors.New("pq: password=secret host=db"), 500, `{"error":"internal server error"}`},
		{"not found", fmt.Errorf("find: %w", domain.ErrNotFound), 404, `{"error":"not found"}`},
		{"email taken", domain.ErrEmailTaken, 409, `{"error":"email already registered"}`},
		{"todo done", domain.ErrTodoDone, 409, `{"error":"completed tasks can't be edited; mark it as not done first"}`},
		{"bad credentials", domain.ErrInvalidCredentials, 401, `{"error":"invalid email or password"}`},
		{"bad token", domain.ErrInvalidToken, 401, `{"error":"invalid or expired token"}`},
		{"fiber error", fiber.NewError(400, "invalid JSON body"), 400, `{"error":"invalid JSON body"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
			app.Get("/", func(fiber.Ctx) error { return tt.err })

			status, body := doRequest(t, app, http.MethodGet, "/", "")
			if status != tt.wantStatus || body != tt.wantBody {
				t.Fatalf("got %d %s, want %d %s", status, body, tt.wantStatus, tt.wantBody)
			}
		})
	}
}

// TestParseBody checks JSON parsing and validation messages
func TestParseBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{"valid", `{"email":"a@example.com","password":"password123"}`, 200, `ok`},
		{"broken json", `{"email":`, 400, `{"error":"invalid JSON body"}`},
		{"missing fields", `{}`, 400, `{"error":"email is required; password is required"}`},
		{"bad email + short password", `{"email":"nope","password":"123"}`, 400,
			`{"error":"email must be a valid email; password must be at least 8 characters"}`},
		{"password too long", `{"email":"a@example.com","password":"` + strings.Repeat("p", 73) + `"}`, 400,
			`{"error":"password must be at most 72 characters"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
			app.Post("/", func(c fiber.Ctx) error {
				var req authRequest
				if err := parseBody(c, &req); err != nil {
					return err
				}
				return c.SendString("ok")
			})

			status, body := doRequest(t, app, http.MethodPost, "/", tt.body)
			if status != tt.wantStatus || body != tt.wantBody {
				t.Fatalf("got %d %s, want %d %s", status, body, tt.wantStatus, tt.wantBody)
			}
		})
	}
}
