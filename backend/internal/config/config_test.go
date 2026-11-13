package config

import (
	"strings"
	"testing"
)

// setValidEnv sets every required env var; t.Setenv restores them after the test
func setValidEnv(t *testing.T) {
	t.Setenv("APP_PORT", "8080")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "todo")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "todo")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
}

// TestLoadConfig
func TestLoadConfig(t *testing.T) {
	setValidEnv(t)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8080" || cfg.DB.Host != "localhost" || cfg.DB.Password != "secret" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

// TestLoadConfigJWTSecret checks that a missing or short secret stops startup
func TestLoadConfigJWTSecret(t *testing.T) {
	for _, secret := range []string{"", strings.Repeat("x", 31)} {
		setValidEnv(t)
		t.Setenv("JWT_SECRET", secret)

		if _, err := LoadConfig(); err == nil {
			t.Fatalf("expected error for secret of length %d", len(secret))
		}
	}
}

// TestDSN checks that special chars in the password are escaped
func TestDSN(t *testing.T) {
	d := DB{Host: "db", Port: "5432", User: "todo", Password: "p@ss word/#", Name: "todo", SSLMode: "disable"}

	got := d.DSN()
	want := "postgres://todo:p%40ss%20word%2F%23@db:5432/todo?sslmode=disable"
	if got != want {
		t.Fatalf("DSN = %s, want %s", got, want)
	}
	if strings.Contains(got, "p@ss") {
		t.Fatal("password was not escaped")
	}
}
