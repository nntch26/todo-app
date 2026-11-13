package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nntch26/todo-app/backend/internal/domain"
)

var testSecret = strings.Repeat("s", 32)

// fakeUserStore is an in-memory UserStore, so tests don't need a database
type fakeUserStore struct {
	byEmail map[string]domain.User
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{byEmail: map[string]domain.User{}}
}

func (f *fakeUserStore) Create(_ context.Context, u *domain.User) error {
	if _, ok := f.byEmail[u.Email]; ok {
		return domain.ErrEmailTaken
	}
	u.ID = "user-" + u.Email
	f.byEmail[u.Email] = *u
	return nil
}

func (f *fakeUserStore) FindByEmail(_ context.Context, email string) (domain.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (f *fakeUserStore) FindByID(_ context.Context, id string) (domain.User, error) {
	for _, u := range f.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

// newTestService returns a service with one registered user: a@example.com / password123
func newTestService(t *testing.T) *AuthService {
	t.Helper()
	s := NewAuthService(newFakeUserStore(), testSecret)
	if _, err := s.Register(context.Background(), "a@example.com", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}
	return s
}

// TestRegister checks hashing and email normalization
func TestRegister(t *testing.T) {
	s := NewAuthService(newFakeUserStore(), testSecret)

	u, err := s.Register(context.Background(), "  A@Example.COM ", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "a@example.com" {
		t.Errorf("email = %q, want normalized", u.Email)
	}
	if u.PasswordHash == "" || u.PasswordHash == "password123" {
		t.Error("password must be stored as a hash")
	}
}

// TestRegisterDuplicate checks that the same email (any case) can't register twice
func TestRegisterDuplicate(t *testing.T) {
	s := newTestService(t)

	_, err := s.Register(context.Background(), "A@EXAMPLE.com", "password456")
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("err = %v, want ErrEmailTaken", err)
	}
}

// TestLogin covers success and both failure cases
func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{name: "correct", email: "a@example.com", password: "password123"},
		{name: "email case-insensitive", email: "A@Example.com", password: "password123"},
		{name: "wrong password", email: "a@example.com", password: "wrong-pass", wantErr: domain.ErrInvalidCredentials},
		{name: "unknown email", email: "nobody@example.com", password: "password123", wantErr: domain.ErrInvalidCredentials},
	}

	s := newTestService(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, token, err := s.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Token must round-trip back to the same user
			id, err := s.ParseToken(token)
			if err != nil || id != u.ID {
				t.Fatalf("ParseToken = %q, %v; want %q", id, err, u.ID)
			}
		})
	}
}

// TestParseTokenRejects checks tokens that must not be accepted
func TestParseTokenRejects(t *testing.T) {
	s := newTestService(t)
	_, valid, err := s.Login(context.Background(), "a@example.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	// Signed with a different secret
	other := NewAuthService(newFakeUserStore(), strings.Repeat("o", 32))
	otherToken, _ := other.issueToken("user-a@example.com")

	// alg "none" = unsigned token
	noneToken, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject:   "user-a@example.com",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	// Issued 25h ago, so already expired
	past := NewAuthService(newFakeUserStore(), testSecret)
	past.now = func() time.Time { return time.Now().Add(-25 * time.Hour) }
	expiredToken, _ := past.issueToken("user-a@example.com")

	tests := map[string]string{
		"empty":         "",
		"garbage":       "not-a-jwt",
		"tampered":      valid[:len(valid)-2] + "xx",
		"different key": otherToken,
		"alg none":      noneToken,
		"expired":       expiredToken,
	}
	for name, token := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := s.ParseToken(token); !errors.Is(err, domain.ErrInvalidToken) {
				t.Fatalf("err = %v, want ErrInvalidToken", err)
			}
		})
	}
}
