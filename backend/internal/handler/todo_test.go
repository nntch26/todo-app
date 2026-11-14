package handler

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// TestUpdateRequestDueDate: missing key, null and a value must be three different things
func TestUpdateRequestDueDate(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantSet bool
		wantNil bool
	}{
		{"missing = keep", `{"done":true}`, false, true},
		{"null = remove", `{"due_date":null}`, true, true},
		{"value = change", `{"due_date":"2026-12-31"}`, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req updateTodoRequest
			err := json.Unmarshal([]byte(tt.body), &req)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			gotNil := req.DueDate.Value == nil
			if req.DueDate.Set != tt.wantSet || gotNil != tt.wantNil {
				t.Errorf("got Set=%v Value=%v", req.DueDate.Set, req.DueDate.Value)
			}
		})
	}
}

// TestParseDueDate: nil, a valid date, and bad formats
func TestParseDueDate(t *testing.T) {
	// nil -> no due date, no error
	date, err := parseDueDate(nil)
	if date != nil || err != nil {
		t.Errorf("nil: got %v, %v", date, err)
	}

	// valid date
	valid := "2026-12-31"
	date, err = parseDueDate(&valid)
	if err != nil || date.Format(dateLayout) != valid {
		t.Errorf("valid: got %v, %v", date, err)
	}

	// bad formats -> error
	badInputs := []string{"", "31/12/2026", "2026-13-01", "2026-12-31T00:00:00Z"}
	for _, bad := range badInputs {
		_, err := parseDueDate(&bad)
		if err == nil {
			t.Errorf("%q: want error", bad)
		}
	}
}

// TestTodoIDNotUUID: a bad id answers 404 before reaching the DB
func TestTodoIDNotUUID(t *testing.T) {
	// Small app with one route that only calls todoID
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Delete("/todos/:id", func(c fiber.Ctx) error {
		id, err := todoID(c)
		if err != nil {
			return err
		}
		return c.SendString(id)
	})

	status, body := doRequest(t, app, "DELETE", "/todos/abc", "")
	if status != 404 || body != `{"error":"not found"}` {
		t.Errorf("non-uuid: got %d %s", status, body)
	}

	validID := "6f1c1a52-6a43-4b6e-9a3e-2d1f4c7b8e90"
	status, body = doRequest(t, app, "DELETE", "/todos/"+validID, "")
	if status != 200 || body != validID {
		t.Errorf("uuid: got %d %s", status, body)
	}
}

// TestUpdateRequestValidation: empty title is rejected, missing title is fine
func TestUpdateRequestValidation(t *testing.T) {
	// Small app with one route that only calls parseBody
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Patch("/", func(c fiber.Ctx) error {
		var req updateTodoRequest
		err := parseBody(c, &req)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	status, body := doRequest(t, app, "PATCH", "/", `{"title":""}`)
	if status != 400 || body != `{"error":"title must not be empty"}` {
		t.Errorf("empty title: got %d %s", status, body)
	}

	status, _ = doRequest(t, app, "PATCH", "/", `{"done":true}`)
	if status != 204 {
		t.Errorf("only done: got %d", status)
	}
}
