package commands

import (
	"context"
	"fmt"

	"github.com/nedson202/go-cqrs/internal/domain/commands"
	"github.com/nedson202/go-cqrs/internal/domain/events"
)

// UserCommandHandler handles user-related commands
type UserCommandHandler struct {
	eventStore events.EventStore
}

// NewUserCommandHandler creates a new UserCommandHandler
func NewUserCommandHandler(store events.EventStore) *UserCommandHandler {
	return &UserCommandHandler{
		eventStore: store,
	}
}

// Handle implements the CommandHandler interface
func (h *UserCommandHandler) Handle(ctx context.Context, cmd commands.Command) error {
	switch c := cmd.(type) {
	case *commands.CreateUser:
		return h.handleCreateUser(ctx, c)
	case *commands.UpdateUser:
		return h.handleUpdateUser(ctx, c)
	default:
		return fmt.Errorf("unknown command type: %T", cmd)
	}
}


// Ensure UserCommandHandler implements CommandHandler interface
var _ commands.CommandHandler = (*UserCommandHandler)(nil) 
