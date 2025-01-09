package commands

// UpdateUser command
type UpdateUser struct {
    ID       string `json:"id"`
    Username string `json:"username"`
}

func (c *UpdateUser) Type() string {
    return "UpdateUser"
}

func (c *UpdateUser) Validate() error {
    if c.ID == "" {
        return ErrIDRequired
    }
    if c.Username == "" {
        return ErrUsernameRequired
    }
    return nil
}

// Ensure command implements the Command interface
var _ Command = (*UpdateUser)(nil) 
