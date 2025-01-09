package events

import (
	"encoding/json"

	"github.com/nedson202/go-cqrs/internal/domain/types"
)

// UserCreatedEvent represents the event when a user is created
type UserCreatedEvent struct {
	ID       types.ID
	Email    string
	Username string
}

func NewUserCreatedEvent(id types.ID, email, username string) (*Event, error) {
	data, err := json.Marshal(UserCreatedEvent{
		ID:       id,
		Email:    email,
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return &Event{
		EventID:       id.String(),
		AggregateType: "user",
		EventType:     "UserCreated",
		Data:         data,
	}, nil
}

// UserUpdatedEvent represents the event when a user is updated
type UserUpdatedEvent struct {
	ID       types.ID
	Username string
}

func NewUserUpdatedEvent(id types.ID, username string, version int) (*Event, error) {
	data, err := json.Marshal(UserUpdatedEvent{
		ID:       id,
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return &Event{
		EventID:       id.String(),
		AggregateType: "user",
		EventType:     "UserUpdated",
		Version:       version,
		Data:         data,
	}, nil
} 
