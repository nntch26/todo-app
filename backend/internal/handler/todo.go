package handler

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nntch26/todo-app/backend/internal/domain"
	"github.com/nntch26/todo-app/backend/internal/middleware"
	"github.com/nntch26/todo-app/backend/internal/service"
)

// dateLayout is the due_date format in JSON: "2026-12-31"
// (Go writes date formats using the example date 2006-01-02)
const dateLayout = "2006-01-02"

// TodoHandler handles HTTP for /api/todos; the rules live in service.TodoService
type TodoHandler struct {
	todos *service.TodoService
}

func NewTodoHandler(todos *service.TodoService) *TodoHandler {
	return &TodoHandler{todos: todos}
}

// ---------- Request / response bodies ----------

type createTodoRequest struct {
	Title       string  `json:"title" validate:"required,max=200"`
	Description string  `json:"description" validate:"max=2000"`
	DueDate     *string `json:"due_date"` // null or missing = no due date
}

type updateTodoRequest struct {
	Title       *string      `json:"title" validate:"omitnil,min=1,max=200"`
	Description *string      `json:"description" validate:"omitnil,max=2000"`
	Done        *bool        `json:"done"`
	DueDate     optionalDate `json:"due_date"`
}

//
type optionalDate struct {
	Set   bool
	Value *string
}

// UnmarshalJSON is called by encoding/json only when "due_date" is in the body
func (o *optionalDate) UnmarshalJSON(b []byte) error {
	o.Set = true
	return json.Unmarshal(b, &o.Value)
}

// todoResponse is the JSON we send back for one todo
type todoResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     *string   `json:"due_date"` // "2026-12-31" or null
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// toTodoResponse converts a domain.Todo into the JSON shape
func toTodoResponse(t domain.Todo) todoResponse {
	// time.Time -> "2026-12-31" (only if there is a due date)
	var dueDate *string
	if t.DueDate != nil {
		formatted := t.DueDate.Format(dateLayout)
		dueDate = &formatted
	}

	return todoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     dueDate,
		Done:        t.Done,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// ---------- Handlers (all routes are behind RequireAuth) ----------

// List GET /api/todos
func (h *TodoHandler) List(c fiber.Ctx) error {
	// Who is calling (set by middleware.RequireAuth)
	userID := middleware.UserID(c)

	// Call service -> repository -> DB
	todos, err := h.todos.List(c.Context(), userID)
	if err != nil {
		return err
	}

	// Convert each todo to JSON shape
	response := make([]todoResponse, 0, len(todos))
	for _, t := range todos {
		response = append(response, toTodoResponse(t))
	}
	return c.JSON(response)
}

// Create POST /api/todos
func (h *TodoHandler) Create(c fiber.Ctx) error {
	userID := middleware.UserID(c)

	// Read + validate JSON body (400 if invalid)
	var req createTodoRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}

	// "2026-12-31" -> time.Time (400 if wrong format)
	dueDate, err := parseDueDate(req.DueDate)
	if err != nil {
		return err
	}

	newTodo := domain.Todo{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}

	// Call service: sets owner + done=false, then saves
	created, err := h.todos.Create(c.Context(), userID, newTodo)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(toTodoResponse(created))
}

// Update PATCH /api/todos/:id
func (h *TodoHandler) Update(c fiber.Ctx) error {
	userID := middleware.UserID(c)

	id, err := todoID(c)
	if err != nil {
		return err
	}

	var req updateTodoRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}

	patch := domain.TodoPatch{
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}

	// due_date only 
	if req.DueDate.Set {
		dueDate, err := parseDueDate(req.DueDate.Value)
		if err != nil {
			return err
		}
		patch.DueDate = &domain.DueDateChange{Date: dueDate}
	}

	// Call service: find own todo -> apply patch -> save (404 if not this user's)
	updated, err := h.todos.Update(c.Context(), userID, id, patch)
	if err != nil {
		return err
	}
	return c.JSON(toTodoResponse(updated))
}

// Delete DELETE /api/todos/:id
func (h *TodoHandler) Delete(c fiber.Ctx) error {
	userID := middleware.UserID(c)

	// Read :id from the URL (404 if not a UUID)
	id, err := todoID(c)
	if err != nil {
		return err
	}

	// Call service
	err = h.todos.Delete(c.Context(), userID, id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ---------- Helpers ----------

// todoID reads :id from the URL.
func todoID(c fiber.Ctx) (string, error) {
	id := c.Params("id")

	err := validate.Var(id, "uuid")
	if err != nil {
		return "", domain.ErrNotFound
	}
	return id, nil
}

// parseDueDate turns "2026-12-31" into a time.Time
func parseDueDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	date, err := time.Parse(dateLayout, *value)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "due_date must be a date like 2026-12-31")
	}
	return &date, nil
}
