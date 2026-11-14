package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// TodoRepository reads/writes the todos table with GORM.
// Every query is filtered by user_id, so a user can never touch another user's todo.
type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// ListByUser returns the user's todos, newest first
func (r *TodoRepository) ListByUser(ctx context.Context, userID string) ([]domain.Todo, error) {
	todos := []domain.Todo{} // empty slice, so JSON is [] not null

	// SELECT * FROM todos WHERE user_id = $1 ORDER BY created_at DESC
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&todos).Error

	return todos, err
}

// Create inserts a todo; the DB fills in ID, CreatedAt, UpdatedAt
func (r *TodoRepository) Create(ctx context.Context, t *domain.Todo) error {
	// INSERT INTO todos (...) VALUES (...)
	// Omit associations: don't let GORM also try to insert t.User
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Create(t).Error
}

// FindByID returns one todo, only if it belongs to userID
func (r *TodoRepository) FindByID(ctx context.Context, userID, id string) (domain.Todo, error) {
	var t domain.Todo

	// SELECT * FROM todos WHERE id = $1 AND user_id = $2 LIMIT 1
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&t).Error

	// gorm.ErrRecordNotFound -> domain.ErrNotFound
	return t, translateNotFound(err)
}

// Update saves the editable fields of t
func (r *TodoRepository) Update(ctx context.Context, t *domain.Todo) error {
	// Use a map, not the struct: with a struct GORM skips zero values,
	// so done=false or description="" would never be saved
	changes := map[string]any{
		"title":       t.Title,
		"description": t.Description,
		"due_date":    t.DueDate,
		"done":        t.Done,
	}

	// UPDATE todos SET ..., updated_at = now WHERE id = $1 AND user_id = $2
	result := r.db.WithContext(ctx).
		Model(&domain.Todo{}).
		Where("id = ? AND user_id = ?", t.ID, t.UserID).
		Updates(changes)

	return rowsOrNotFound(result)
}

// Delete removes one todo, only if it belongs to userID
func (r *TodoRepository) Delete(ctx context.Context, userID, id string) error {
	// DELETE FROM todos WHERE id = $1 AND user_id = $2
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Todo{})

	return rowsOrNotFound(result)
}

// rowsOrNotFound: 0 rows changed = the todo doesn't exist (or isn't this user's)
func rowsOrNotFound(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
