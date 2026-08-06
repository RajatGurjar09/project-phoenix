package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	hostEnv     = "PHOENIX_DB_HOST"
	portEnv     = "PHOENIX_DB_PORT"
	nameEnv     = "PHOENIX_DB_NAME"
	userEnv     = "PHOENIX_DB_USER"
	passwordEnv = "PHOENIX_DB_PASSWORD"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadConfig() Config {
	return Config{
		Host:     envOrDefault(hostEnv, "postgres"),
		Port:     envOrDefault(portEnv, "5432"),
		Name:     envOrDefault(nameEnv, "phoenix"),
		User:     envOrDefault(userEnv, "phoenix"),
		Password: os.Getenv(passwordEnv),
	}
}

func (c Config) ConnectionString() string {
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   "/" + c.Name,
	}

	return connectionURL.String()
}

func Connect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func envOrDefault(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
