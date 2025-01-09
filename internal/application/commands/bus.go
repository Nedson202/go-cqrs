package commands

import (
	"context"
	"sync"

	"github.com/nedson202/go-cqrs/internal/domain/commands"
	"github.com/nedson202/go-cqrs/internal/pkg/errors"
)

var (
	ErrHandlerNotFound = errors.NewCommandError(
		"HANDLER_NOT_FOUND",
		"no handler registered for command",
		nil,
	)
)

// CommandBus manages command handlers and dispatches commands
type CommandBus struct {
	handlers map[string]commands.CommandHandler
	mu       sync.RWMutex
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]commands.CommandHandler),
	}
}

func (b *CommandBus) Register(cmdType string, handler commands.CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[cmdType] = handler
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd commands.Command) error {
	b.mu.RLock()
	handler, exists := b.handlers[cmd.Type()]
	b.mu.RUnlock()

	if !exists {
		return ErrHandlerNotFound
	}

	return handler.Handle(ctx, cmd)
} 
