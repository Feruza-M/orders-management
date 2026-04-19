package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func ConfigFromEnv() Config {
	return Config{
		Host:     envOrDefault("DB_HOST", "127.0.0.1"),
		Port:     envOrDefault("DB_PORT", "5432"),
		Name:     envOrDefault("DB_NAME", "orders_db"),
		User:     envOrDefault("DB_USER", "orders_app"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func NewPool(cfg Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pgxpool.New(ctx, cfg.DSN())
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
