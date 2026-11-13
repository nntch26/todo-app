package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a user; DB fills in ID and CreatedAt
func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrEmailTaken
	}
	return err
}

// FindByEmail
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, translateNotFound(err)
}

// FindByID
func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, translateNotFound(err)
}

// translateNotFound hides GORM's error type from the service layer
func translateNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}
