package validation

// CommandValidator defines the interface for command validation
type CommandValidator interface {
	Validate(command interface{}) error
}
