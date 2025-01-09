package queries

import (
	"context"
	"fmt"
	"sync"
)

// Query represents a query that can be handled
type Query interface {
	Type() string
}

// QueryHandler defines how queries should be handled
type QueryHandler interface {
	Handle(context.Context, Query) (interface{}, error)
}

// QueryBus manages query dispatch
type QueryBus struct {
	handlers map[string]QueryHandler
	mu       sync.RWMutex
}

func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

func (b *QueryBus) Register(queryType string, handler QueryHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[queryType] = handler
}

func (b *QueryBus) Ask(ctx context.Context, q Query) (interface{}, error) {
	b.mu.RLock()
	handler, exists := b.handlers[q.Type()]
	b.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for query type: %s", q.Type())
	}

	return handler.Handle(ctx, q)
} 
