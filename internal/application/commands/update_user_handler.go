package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainCommands "github.com/nedson202/go-cqrs/internal/domain/commands"
	"github.com/nedson202/go-cqrs/internal/domain/events"
)

func (h *UserCommandHandler) handleUpdateUser(ctx context.Context, cmd *domainCommands.UpdateUser) error {
    if err := cmd.Validate(); err != nil {
        return fmt.Errorf("validate command: %w", err)
    }

    userID, err := uuid.Parse(cmd.ID)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }

    // Get current version
    userEvents, err := h.eventStore.Load(ctx, userID.String(), 0)
    if err != nil {
        return fmt.Errorf("load events: %w", err)
    }

    currentVersion := len(userEvents)
    event, err := events.NewUserUpdatedEvent(userID, cmd.Username, currentVersion)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	if err := h.eventStore.Save(ctx, []events.Event{*event}); err != nil {
		return fmt.Errorf("save event: %w", err)
	}

	return nil
}
