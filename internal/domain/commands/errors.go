package commands

import "errors"

var (
    ErrEmailRequired    = errors.New("email is required")
    ErrUsernameRequired = errors.New("username is required")
    ErrIDRequired       = errors.New("id is required")
) 
