package commands

import "context"

// Command represents a domain command
type Command interface {
    // Type returns the type of the command
    Type() string
    // Validate validates the command
    Validate() error
}

// CommandHandler defines how commands should be handled
type CommandHandler interface {
    // Handle processes a command
    Handle(ctx context.Context, cmd Command) error
} 
