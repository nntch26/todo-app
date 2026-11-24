package service

import (
	"context"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// TodoStore is what TodoService needs from the data layer.
// Real app: repository.TodoRepository (PostgreSQL). Tests: fakeTodoStore (in memory).
// Every method takes userID: a todo owned by someone else is reported as domain.ErrNotFound.
type TodoStore interface {
	ListByUser(ctx context.Context, userID string) ([]domain.Todo, error)
	Create(ctx context.Context, t *domain.Todo) error
	FindByID(ctx context.Context, userID, id string) (domain.Todo, error)
	Update(ctx context.Context, t *domain.Todo) error
	Delete(ctx context.Context, userID, id string) error
}

// TodoService holds the todo rules; it never talks to the DB directly, only through TodoStore
type TodoService struct {
	todos TodoStore
}

func NewTodoService(todos TodoStore) *TodoService {
	return &TodoService{todos: todos}
}

// List returns the user's own todos
func (s *TodoService) List(ctx context.Context, userID string) ([]domain.Todo, error) {
	// Call repository: SELECT ... WHERE user_id = ?
	return s.todos.ListByUser(ctx, userID)
}

// Create saves a new todo for userID
func (s *TodoService) Create(ctx context.Context, userID string, t domain.Todo) (domain.Todo, error) {
	// Rules: the server decides id, owner and done; never trust these from the client
	t.ID = "" // empty -> DB generates a UUID
	t.UserID = userID
	t.Done = false

	// Call repository: INSERT (fills t.ID, CreatedAt, UpdatedAt)
	err := s.todos.Create(ctx, &t)
	if err != nil {
		return domain.Todo{}, err
	}
	return t, nil
}

// Update changes only the fields set in the patch.
// Sending the same patch twice gives the same result (idempotent).
func (s *TodoService) Update(ctx context.Context, userID, id string, p domain.TodoPatch) (domain.Todo, error) {
	// 1. Call repository: load the todo (ErrNotFound if missing or not this user's)
	t, err := s.todos.FindByID(ctx, userID, id)
	if err != nil {
		return domain.Todo{}, err
	}

	// 2. Rule: a completed todo is read-only; only "done" may change (uncheck first, then edit)
	editsContent := p.Title != nil || p.Description != nil || p.DueDate != nil
	if t.Done && editsContent {
		return domain.Todo{}, domain.ErrTodoDone
	}

	// 3. Apply the patch in memory (no DB yet); nil = not sent = keep
	if p.Title != nil {
		t.Title = *p.Title
	}
	if p.Description != nil {
		t.Description = *p.Description
	}
	if p.Done != nil {
		t.Done = *p.Done
	}
	if p.DueDate != nil {
		t.DueDate = p.DueDate.Date // nil Date = remove the due date
	}

	// 4. Call repository: save the changes (UPDATE ... WHERE id = ? AND user_id = ?)
	err = s.todos.Update(ctx, &t)
	if err != nil {
		return domain.Todo{}, err
	}

	// 5. Call repository again: read back to get the new updated_at
	return s.todos.FindByID(ctx, userID, id)
}

// Delete removes the user's own todo (ErrNotFound if missing or not this user's)
func (s *TodoService) Delete(ctx context.Context, userID, id string) error {
	// Call repository: DELETE ... WHERE id = ? AND user_id = ?
	return s.todos.Delete(ctx, userID, id)
}
