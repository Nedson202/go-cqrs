package identity

import (
	"github.com/google/uuid"
	"github.com/nedson202/go-cqrs/internal/domain/types"
)

type UUID struct {
	id uuid.UUID
}

func NewUUID() *UUID {
	return &UUID{id: uuid.New()}
}

func (u *UUID) String() string {
	return u.id.String()
}

// Ensure UUID implements ID interface
var _ types.ID = (*UUID)(nil) 
