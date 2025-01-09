package events

import (
	"context"
)

// Event represents a domain event
type Event struct {
	EventID       string
	AggregateID   string
	AggregateType string
	EventType     string
	Version       int
	Data          []byte
	Metadata      []byte
	Timestamp     int64
}

// EventStore defines the domain interface for event storage
type EventStore interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
	Stream(ctx context.Context, aggregateID string, fromVersion int) (<-chan Event, <-chan error)
}

// Ensure EventStore implements EventStore interface
var _ EventStore = (EventStore)(nil) 
