package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// TestHealth checks that GET /api/health returns 200 and {"status":"ok"}.
func TestHealth(t *testing.T) {
	app := fiber.New()
	app.Get("/api/health", Health)

	// Send fake request (no real server/port needed)
	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if got, want := string(body), `{"status":"ok"}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
