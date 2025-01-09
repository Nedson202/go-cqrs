package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainCommands "github.com/nedson202/go-cqrs/internal/domain/commands"
	"github.com/nedson202/go-cqrs/internal/domain/events"
)

func (h *UserCommandHandler) handleCreateUser(ctx context.Context, cmd *domainCommands.CreateUser) error {
    if err := cmd.Validate(); err != nil {
        return fmt.Errorf("validate command: %w", err)
    }

    userID := uuid.New()
    event, err := events.NewUserCreatedEvent(userID, cmd.Email, cmd.Username)
    if err != nil {
        return fmt.Errorf("create event: %w", err)
    }

    if err := h.eventStore.Save(ctx, []events.Event{*event}); err != nil {
        return fmt.Errorf("save event: %w", err)
    }

    return nil
} 
