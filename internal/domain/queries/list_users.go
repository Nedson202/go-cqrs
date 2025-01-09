package queries

// ListUsers query
type ListUsers struct{}

func (q *ListUsers) Type() string {
    return "ListUsers"
}

// Ensure query implements the Query interface
var _ Query = (*ListUsers)(nil) 