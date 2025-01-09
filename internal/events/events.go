package events

import "time"

type UserEvent struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    Data      []byte    `json:"data"`
    Timestamp time.Time `json:"timestamp"`
}

type EventStore interface {
    Save(event UserEvent) error
    Get(id string) ([]UserEvent, error)
} 
