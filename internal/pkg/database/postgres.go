package database

import (
	"context"
	"fmt"

	"asona/config"
	"asona/internal/repository/db/pg"
)

// NewPostgresConnection creates a new PostgreSQL connection and returns a BeginnerExecutor
func NewPostgresConnection() (pg.BeginnerExecutor, error) {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	return pg.NewPool(
		dsn,
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
	)
}

// CheckConnection verifies database connectivity
func CheckConnection(db pg.BeginnerExecutor) error {
	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("database connection check failed: %w", err)
	}
	return nil
}
