package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"asona/internal/pkg/database"
	"asona/internal/repository/db/pg"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error

	// DB returns the underlying pg.BeginnerExecutor.
	DB() pg.BeginnerExecutor
}

type service struct {
	db pg.BeginnerExecutor
}

var (
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to create postgres connection: %v", err)
	}

	// Verify connectivity on startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	dbInstance = &service{db: db}
	return dbInstance
}

// DB returns the underlying pg.BeginnerExecutor for use in repositories.
func (s *service) DB() pg.BeginnerExecutor {
	return s.db
}

// Health checks the health of the Postgres connection.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	if err := s.db.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("Postgres health check failed: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

// Close closes the Postgres connection.
func (s *service) Close() error {
	log.Printf("Disconnected from Postgres")
	return s.db.Close()
}
