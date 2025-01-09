package validation

import (
	"fmt"

	domainCommands "github.com/nedson202/go-cqrs/internal/domain/commands"
)

// UserCommandValidator handles validation for user-related commands
type UserCommandValidator struct{}

func NewUserCommandValidator() *UserCommandValidator {
    return &UserCommandValidator{}
}

func (v *UserCommandValidator) Validate(command interface{}) error {
    switch cmd := command.(type) {
    case *domainCommands.CreateUser:
        return v.validateCreateUser(cmd)
    case *domainCommands.UpdateUser:
        return v.validateUpdateUser(cmd)
    default:
        return fmt.Errorf("unknown command type: %T", command)
    }
}

func (v *UserCommandValidator) validateCreateUser(cmd *domainCommands.CreateUser) error {
    return cmd.Validate()
}

func (v *UserCommandValidator) validateUpdateUser(cmd *domainCommands.UpdateUser) error {
    return cmd.Validate()
}

// Ensure UserCommandValidator implements CommandValidator
var _ CommandValidator = (*UserCommandValidator)(nil) 
