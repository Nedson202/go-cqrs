package queries

import (
	"context"
	"fmt"

	domainQueries "github.com/nedson202/go-cqrs/internal/domain/queries"
)

func (h *UserQueryHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	switch query := q.(type) {
    case *domainQueries.GetUser:
        return h.handleGetUser(ctx, query)
    case *domainQueries.ListUsers:
        return h.handleListUsers(ctx)
    default:
        return nil, fmt.Errorf("unknown query type: %T", q)
    }
}

