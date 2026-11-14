package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// ---------- Fake store (in memory, replaces the real DB in tests) ----------

// fakeTodoStore implements TodoStore with a map, scoped by user like the real repository
type fakeTodoStore struct {
	todos  map[string]domain.Todo
	nextID int
}

func newFakeTodoStore() *fakeTodoStore {
	return &fakeTodoStore{todos: map[string]domain.Todo{}}
}

// ListByUser
func (f *fakeTodoStore) ListByUser(_ context.Context, userID string) ([]domain.Todo, error) {
	result := []domain.Todo{}
	for _, t := range f.todos {
		if t.UserID == userID {
			result = append(result, t)
		}
	}
	return result, nil
}

// Create
func (f *fakeTodoStore) Create(_ context.Context, t *domain.Todo) error {
	f.nextID++
	t.ID = fmt.Sprintf("todo-%d", f.nextID)
	f.todos[t.ID] = *t
	return nil
}

// FindByID
func (f *fakeTodoStore) FindByID(_ context.Context, userID, id string) (domain.Todo, error) {
	t, found := f.todos[id]
	if !found || t.UserID != userID {
		return domain.Todo{}, domain.ErrNotFound
	}
	return t, nil
}

// Update
func (f *fakeTodoStore) Update(ctx context.Context, t *domain.Todo) error {
	_, err := f.FindByID(ctx, t.UserID, t.ID)
	if err != nil {
		return err
	}
	f.todos[t.ID] = *t
	return nil
}

// Delete
func (f *fakeTodoStore) Delete(ctx context.Context, userID, id string) error {
	_, err := f.FindByID(ctx, userID, id)
	if err != nil {
		return err
	}
	delete(f.todos, id)
	return nil
}

// ---------- Setup ----------

// newTodoTest creates a service where alice owns one todo with a due date
func newTodoTest(t *testing.T) (*TodoService, domain.Todo) {
	t.Helper()
	s := NewTodoService(newFakeTodoStore())

	dueDate := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	aliceTodo := domain.Todo{
		Title:       "buy milk",
		Description: "2 bottles",
		DueDate:     &dueDate,
	}

	created, err := s.Create(context.Background(), "alice", aliceTodo)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	return s, created
}

// ---------- Tests ----------

// TestCreateTodo: id, owner and done sent by the client must be ignored
func TestCreateTodo(t *testing.T) {
	s := NewTodoService(newFakeTodoStore())

	fromClient := domain.Todo{ID: "hacked", UserID: "bob", Title: "x", Done: true}
	got, err := s.Create(context.Background(), "alice", fromClient)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if got.ID == "hacked" || got.UserID != "alice" || got.Done {
		t.Errorf("got %+v, want new ID, owner alice, not done", got)
	}
}

// TestListTodosOnlyOwn: each user sees only their own todos
func TestListTodosOnlyOwn(t *testing.T) {
	s, _ := newTodoTest(t)
	ctx := context.Background()

	aliceList, _ := s.List(ctx, "alice")
	if len(aliceList) != 1 {
		t.Errorf("alice sees %d todos, want 1", len(aliceList))
	}

	bobList, _ := s.List(ctx, "bob")
	if len(bobList) != 0 {
		t.Errorf("bob sees %d todos, want 0", len(bobList))
	}
}

// TestUpdateTodo
func TestUpdateTodo(t *testing.T) {
	ctx := context.Background()

	t.Run("only patched fields change", func(t *testing.T) {
		s, todo := newTodoTest(t)

		newTitle := "buy bread"
		got, err := s.Update(ctx, "alice", todo.ID, domain.TodoPatch{Title: &newTitle})
		if err != nil {
			t.Fatalf("Update: %v", err)
		}

		if got.Title != "buy bread" || got.Description != "2 bottles" || got.DueDate == nil {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("mark done twice gives the same result", func(t *testing.T) {
		s, todo := newTodoTest(t)

		done := true
		patch := domain.TodoPatch{Done: &done}
		first, _ := s.Update(ctx, "alice", todo.ID, patch)
		second, err := s.Update(ctx, "alice", todo.ID, patch)

		if err != nil || !first.Done || !second.Done {
			t.Errorf("first=%v second=%v err=%v, want both done", first.Done, second.Done, err)
		}
	})

	t.Run("clear due date", func(t *testing.T) {
		s, todo := newTodoTest(t)

		removeDueDate := domain.TodoPatch{DueDate: &domain.DueDateChange{Date: nil}}
		got, _ := s.Update(ctx, "alice", todo.ID, removeDueDate)

		if got.DueDate != nil {
			t.Errorf("due date = %v, want nil", got.DueDate)
		}
	})

	t.Run("someone else's todo is not found", func(t *testing.T) {
		s, todo := newTodoTest(t)

		done := true
		_, err := s.Update(ctx, "bob", todo.ID, domain.TodoPatch{Done: &done})
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("err = %v, want ErrNotFound", err)
		}

		// alice's todo must be unchanged
		aliceList, _ := s.List(ctx, "alice")
		if aliceList[0].Done {
			t.Error("bob changed alice's todo")
		}
	})
}

// TestDeleteTodo: only the owner can delete, and only once
func TestDeleteTodo(t *testing.T) {
	ctx := context.Background()
	s, todo := newTodoTest(t)

	err := s.Delete(ctx, "bob", todo.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("bob delete: err = %v, want ErrNotFound", err)
	}

	err = s.Delete(ctx, "alice", todo.ID)
	if err != nil {
		t.Fatalf("alice delete: %v", err)
	}

	err = s.Delete(ctx, "alice", todo.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("second delete: err = %v, want ErrNotFound", err)
	}
}
