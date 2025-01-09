package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nedson202/go-cqrs/internal/application/commands"
	"github.com/nedson202/go-cqrs/internal/application/queries"
	domainCommands "github.com/nedson202/go-cqrs/internal/domain/commands"
	domainModels "github.com/nedson202/go-cqrs/internal/domain/models"
	domainQueries "github.com/nedson202/go-cqrs/internal/domain/queries"
)

type UserHandler struct {
	commandBus *commands.CommandBus
	queryBus   *queries.QueryBus
	logger     *log.Logger
}

func NewUserHandler(cb *commands.CommandBus, qb *queries.QueryBus, l *log.Logger) *UserHandler {
	return &UserHandler{
		commandBus: cb,
		queryBus:   qb,
		logger:     l,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var cmd domainCommands.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.commandBus.Dispatch(r.Context(), &cmd); err != nil {
		h.logger.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	var cmd domainCommands.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.logger.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cmd.ID = id

	if err := h.commandBus.Dispatch(r.Context(), &cmd); err != nil {
		h.logger.Printf("Error updating user: %v", err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := &domainQueries.GetUser{ID: id}
	result, err := h.queryBus.Ask(r.Context(), query)
	if err != nil {
		h.logger.Printf("Error getting user: %v", err)
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	user, ok := result.(*domainModels.UserDTO)
	if !ok {
		h.logger.Printf("Invalid response type from query bus")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	query := &domainQueries.ListUsers{}
	result, err := h.queryBus.Ask(r.Context(), query)
	if err != nil {
		h.logger.Printf("Error listing users: %v", err)
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}

	users, ok := result.([]*domainModels.UserDTO)
	if !ok {
		h.logger.Printf("Invalid response type from query bus")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
} 
