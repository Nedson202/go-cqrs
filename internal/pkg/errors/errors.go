package errors

import "fmt"

type CommandError struct {
    Code    string
    Message string
    Cause   error
}

func (e *CommandError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *CommandError) Unwrap() error {
    return e.Cause
}

func NewCommandError(code string, message string, cause error) *CommandError {
    return &CommandError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
}

type QueryError struct {
    Code    string
    Message string
    Cause   error
}

func (e *QueryError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *QueryError) Unwrap() error {
    return e.Cause
}

func NewQueryError(code string, message string, cause error) *QueryError {
    return &QueryError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
} 
