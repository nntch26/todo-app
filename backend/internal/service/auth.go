// Package service holds business logic. It must not import Fiber or GORM.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

// TokenTTL is how long a login lasts
const TokenTTL = 24 * time.Hour

// UserStore is what AuthService needs from the data layer (implemented by repository.UserRepository)
type UserStore interface {
	Create(ctx context.Context, u *domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type AuthService struct {
	users  UserStore
	secret []byte
	now    func() time.Time // replaceable in tests
}

func NewAuthService(users UserStore, jwtSecret string) *AuthService {
	return &AuthService{users: users, secret: []byte(jwtSecret), now: time.Now}
}

// Register creates a user with a bcrypt-hashed password
func (s *AuthService) Register(ctx context.Context, email, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	u := domain.User{Email: normalizeEmail(email), PasswordHash: string(hash)}
	if err := s.users.Create(ctx, &u); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

// Login checks credentials and returns the user + a signed JWT
func (s *AuthService) Login(ctx context.Context, email, password string) (domain.User, string, error) {
	u, err := s.users.FindByEmail(ctx, normalizeEmail(email))

	// Same error for "no such email" and "wrong password", so attackers can't tell which emails exist
	if errors.Is(err, domain.ErrNotFound) {
		return domain.User{}, "", domain.ErrInvalidCredentials
	}
	if err != nil {
		return domain.User{}, "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return domain.User{}, "", domain.ErrInvalidCredentials
	}

	token, err := s.issueToken(u.ID)
	if err != nil {
		return domain.User{}, "", err
	}
	return u, token, nil
}

// ParseToken verifies signature + expiry and returns the user ID
func (s *AuthService) ParseToken(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(*jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), // reject "none" / other algs
		jwt.WithExpirationRequired(),
		jwt.WithTimeFunc(s.now),
	)
	if err != nil || claims.Subject == "" {
		return "", domain.ErrInvalidToken
	}
	return claims.Subject, nil
}

// Me returns the logged-in user
func (s *AuthService) Me(ctx context.Context, userID string) (domain.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *AuthService) issueToken(userID string) (string, error) {
	now := s.now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(TokenTTL)),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// normalizeEmail so "A@x.com " and "a@x.com" are the same account
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
