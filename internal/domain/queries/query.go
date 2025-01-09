package queries

import "context"

// Query represents a domain query
type Query interface {
    // Type returns the type of the query
    Type() string
}

// QueryHandler defines how queries should be handled
type QueryHandler interface {
    // Handle processes a query and returns a result
    Handle(ctx context.Context, q Query) (interface{}, error)
} 
