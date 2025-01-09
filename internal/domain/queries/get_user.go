package queries

// GetUser query
type GetUser struct {
    ID string `json:"id"`
}

func (q *GetUser) Type() string {
    return "GetUser"
}

// Ensure query implements the Query interface
var _ Query = (*GetUser)(nil) 