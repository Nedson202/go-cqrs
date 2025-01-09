package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nedson202/go-cqrs/internal/domain/models"
	domainQueries "github.com/nedson202/go-cqrs/internal/domain/queries"
)

func (h *UserQueryHandler) handleGetUser(ctx context.Context, query *domainQueries.GetUser) (*models.UserDTO, error) {
    // First try materialized view
    var dto models.UserDTO
    err := h.db.QueryRowContext(ctx, `
        SELECT state
        FROM aggregate_states
        WHERE aggregate_id = $1 AND aggregate_type = 'user'
    `, query.ID).Scan(&dto)

    if err == nil {
        return &dto, nil
    }

    if err != sql.ErrNoRows {
        return nil, fmt.Errorf("query materialized view: %w", err)
    }

    // Fall back to event reconstruction if not in materialized view
    uid, err := uuid.Parse(query.ID)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID: %w", err)
    }

    events, err := h.loadEvents(ctx, uid)
    if err != nil {
		return nil, err
	}

	return h.reconstructState(events)
}
