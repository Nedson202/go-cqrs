package commands

// CreateUser command
type CreateUser struct {
    Email    string `json:"email"`
    Username string `json:"username"`
}

func (c *CreateUser) Type() string {
    return "CreateUser"
}

func (c *CreateUser) Validate() error {
    if c.Email == "" {
        return ErrEmailRequired
    }
    if c.Username == "" {
        return ErrUsernameRequired
    }
    return nil
}

// Ensure command implements the Command interface
var _ Command = (*CreateUser)(nil) 
