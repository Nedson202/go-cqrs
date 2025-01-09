package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/nedson202/go-cqrs/internal/application/commands"
	"github.com/nedson202/go-cqrs/internal/application/queries"
	"github.com/nedson202/go-cqrs/internal/infrastructure/eventstore/postgres"
	"github.com/nedson202/go-cqrs/internal/interfaces/http/handlers"
	"github.com/nedson202/go-cqrs/internal/interfaces/http/router"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Infrastructure
	db, err := setupDatabase()
	if err != nil {
		logger.Fatalf("Failed to setup database: %v", err)
	}
	eventStore := postgres.NewEventStore(db, 1000)

	// Application
	commandBus := commands.NewCommandBus()
	queryBus := queries.NewQueryBus()

	// Register handlers
	userCommandHandler := commands.NewUserCommandHandler(eventStore)
	userQueryHandler := queries.NewUserQueryHandler(db, eventStore)

	commandBus.Register("CreateUser", userCommandHandler)
	commandBus.Register("UpdateUser", userCommandHandler)

	queryBus.Register("GetUser", userQueryHandler)
	queryBus.Register("ListUsers", userQueryHandler)

	// Interface
	userHandler := handlers.NewUserHandler(commandBus, queryBus, logger)
	
	// Setup router with handlers
	routerHandlers := &router.Handlers{
		User: userHandler,
	}

	// Create router with middleware
	r := router.New(routerHandlers)

	// Start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	logger.Printf("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func setupDatabase() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/cqrs_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
} 
