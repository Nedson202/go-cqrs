package validation

// ValidatorFactory creates appropriate validators
type ValidatorFactory struct{}

func NewValidatorFactory() *ValidatorFactory {
    return &ValidatorFactory{}
}

func (f *ValidatorFactory) CreateUserCommandValidator() CommandValidator {
    return NewUserCommandValidator()
}
