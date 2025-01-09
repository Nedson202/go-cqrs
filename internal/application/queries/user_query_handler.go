package queries

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nedson202/go-cqrs/internal/domain/events"
	"github.com/nedson202/go-cqrs/internal/domain/models"
)

// UserQueryHandler handles user-related queries
type UserQueryHandler struct {
    db         *sql.DB
    eventStore events.EventStore
}

func NewUserQueryHandler(db *sql.DB, store events.EventStore) *UserQueryHandler {
    return &UserQueryHandler{
        db:         db,
        eventStore: store,
    }
}
func (h *UserQueryHandler) loadEvents(ctx context.Context, aggregateID uuid.UUID) ([]events.Event, error) {
    var eventList []events.Event
    
    eventCh, errCh := h.eventStore.Stream(ctx, aggregateID.String(), 0)
    
    for {
        select {
        case event, ok := <-eventCh:
            if !ok {
                return eventList, nil
            }
            eventList = append(eventList, event)
        case err := <-errCh:
            if err != nil {
                return nil, fmt.Errorf("stream events: %w", err)
            }
            return eventList, nil
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
}

func (h *UserQueryHandler) reconstructState(events []events.Event) (*models.UserDTO, error) {
    if len(events) == 0 {
        return nil, fmt.Errorf("no events found for user")
    }

    dto := &models.UserDTO{}
    
    for _, event := range events {
        switch event.EventType {
        case "UserCreated":
            var data struct {
                ID       string `json:"id"`
                Email    string `json:"email"`
                Username string `json:"username"`
            }
            if err := json.Unmarshal(event.Data, &data); err != nil {
                return nil, fmt.Errorf("unmarshal UserCreated event: %w", err)
            }
            dto.ID = data.ID
            dto.Email = data.Email
            dto.Username = data.Username

        case "UsernameUpdated":
            var data struct {
                Username string `json:"username"`
            }
            if err := json.Unmarshal(event.Data, &data); err != nil {
                return nil, fmt.Errorf("unmarshal UsernameUpdated event: %w", err)
            }
            dto.Username = data.Username
        }
    }

    return dto, nil
}


func (h *UserQueryHandler) rebuildUsersFromEvents(ctx context.Context) ([]*models.UserDTO, error) {
    rows, err := h.db.QueryContext(ctx, `
        SELECT DISTINCT aggregate_id 
        FROM events 
        WHERE aggregate_type = 'user'
    `)
    if err != nil {
        return nil, fmt.Errorf("query user aggregates: %w", err)
    }
    defer rows.Close()

    var users []*models.UserDTO
    for rows.Next() {
        var idStr string
        if err := rows.Scan(&idStr); err != nil {
            return nil, fmt.Errorf("scan aggregate id: %w", err)
        }

        id, err := uuid.Parse(idStr)
        if err != nil {
            return nil, fmt.Errorf("parse UUID: %w", err)
        }

        events, err := h.loadEvents(ctx, id)
        if err != nil {
            return nil, fmt.Errorf("load events for user %s: %w", id, err)
        }

        user, err := h.reconstructState(events)
        if err != nil {
            return nil, fmt.Errorf("reconstruct user %s: %w", id, err)
        }

        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterate aggregates: %w", err)
    }

    // Refresh the materialized view
    if err := h.refreshMaterializedView(ctx); err != nil {
        log.Printf("Failed to refresh materialized view: %v", err)
    }

    return users, nil
}

func (h *UserQueryHandler) refreshMaterializedView(ctx context.Context) error {
    _, err := h.db.ExecContext(ctx, `
        REFRESH MATERIALIZED VIEW CONCURRENTLY aggregate_states
    `)
    return err
} 

// Ensure UserQueryHandler implements QueryHandler interface
var _ QueryHandler = (*UserQueryHandler)(nil) 

