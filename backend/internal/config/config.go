// Package config loads environment variables at startup.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	DB        DB
}

// minJWTSecretLen: HS256 needs a key of at least 256 bits (32 bytes)
const minJWTSecretLen = 32

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoadConfig
func LoadConfig() (Config, error) {

	// Load .env file
	if err := godotenv.Load("../.env"); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		// Server
		Port: getEnvKey("APP_PORT", "8080"),

		// Auth
		JWTSecret: getEnvKey("JWT_SECRET", ""),

		// Database
		DB: DB{
			Host:     getEnvKey("DB_HOST", ""),
			Port:     getEnvKey("DB_PORT", ""),
			User:     getEnvKey("DB_USER", ""),
			Password: getEnvKey("DB_PASSWORD", ""),
			Name:     getEnvKey("DB_NAME", ""),
			SSLMode:  getEnvKey("DB_SSLMODE", "disable"),
		},
	}

	// Empty/short secret = anyone can forge login tokens, so never start without it
	if len(cfg.JWTSecret) < minJWTSecretLen {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLen)
	}

	return cfg, nil
}

// DSN builds a postgres:// URL; url.UserPassword escapes special chars in the password
func (d DB) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     net.JoinHostPort(d.Host, d.Port),
		Path:     d.Name,
		RawQuery: url.Values{"sslmode": {d.SSLMode}}.Encode(),
	}
	return u.String()
}

// getEnvKey returns the env value, or defaultValue if empty
func getEnvKey(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
