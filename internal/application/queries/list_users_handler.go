package queries

import (
	"context"
	"fmt"

	"github.com/nedson202/go-cqrs/internal/domain/models"
)

func (h *UserQueryHandler) handleListUsers(ctx context.Context) ([]*models.UserDTO, error) {
    // Query the materialized view first
    rows, err := h.db.QueryContext(ctx, `
        SELECT state
        FROM aggregate_states
        WHERE aggregate_type = 'user'
        ORDER BY last_updated DESC
    `)
    if err != nil {
        return nil, fmt.Errorf("query materialized view: %w", err)
    }
    defer rows.Close()

    var users []*models.UserDTO
    for rows.Next() {
        var dto models.UserDTO
        if err := rows.Scan(&dto); err != nil {
            return nil, fmt.Errorf("scan user data: %w", err)
        }
        users = append(users, &dto)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterate rows: %w", err)
    }

    // If no users found in materialized view, rebuild from events
    if len(users) == 0 {
        return h.rebuildUsersFromEvents(ctx)
    }

    return users, nil
}
