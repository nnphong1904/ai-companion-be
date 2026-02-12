package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Addr returns the server listen address.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// DatabaseConfig holds Supabase/PostgreSQL connection settings.
// DATABASE_URL takes precedence when set (recommended for Supabase).
// Falls back to individual DB_* params for local development.
type DatabaseConfig struct {
	// URL is the full connection string (e.g. Supabase pooler URL).
	// Use port 6543 for transaction-mode pooling, 5432 for session mode.
	URL string

	// UsePooler indicates the connection goes through PgBouncer (transaction mode).
	// When true, pgx disables prepared statement caching for compatibility.
	// Supabase pooler URLs (port 6543) require this.
	UsePooler bool

	// Fallback individual params (used when URL is empty).
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns the connection string. DATABASE_URL takes precedence.
func (d DatabaseConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			URL:       os.Getenv("DATABASE_URL"),
			UsePooler: getEnvBool("DB_USE_POOLER", true),
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnvInt("DB_PORT", 5432),
			User:      getEnv("DB_USER", "postgres"),
			Password:  getEnv("DB_PASSWORD", "postgres"),
			Name:      getEnv("DB_NAME", "ai_companion"),
			SSLMode:   getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "change-me-in-production"),
			Expiration: getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
