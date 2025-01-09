package validation

// CommandValidator defines domain validation
type CommandValidator interface {
    Validate(command interface{}) error
} 