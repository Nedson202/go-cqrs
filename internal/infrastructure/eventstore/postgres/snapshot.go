package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Snapshot struct {
    SnapshotID    uuid.UUID       `json:"snapshot_id"`
    AggregateID   uuid.UUID       `json:"aggregate_id"`
    AggregateType string          `json:"aggregate_type"`
    Version       int             `json:"version"`
    State         json.RawMessage `json:"state"`
    CreatedAt     time.Time       `json:"created_at"`
}

func (s *PostgresEventStore) SaveSnapshot(ctx context.Context, snapshot Snapshot) error {
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO snapshots (
            snapshot_id, aggregate_id, aggregate_type, 
            version, state, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `,
        snapshot.SnapshotID,
        snapshot.AggregateID,
        snapshot.AggregateType,
        snapshot.Version,
        snapshot.State,
        snapshot.CreatedAt,
    )
    return err
}

func (s *PostgresEventStore) LoadLatestSnapshot(ctx context.Context, aggregateID uuid.UUID) (*Snapshot, error) {
    var snapshot Snapshot
    err := s.db.QueryRowContext(ctx, `
        SELECT snapshot_id, aggregate_id, aggregate_type, 
               version, state, created_at
        FROM snapshots
        WHERE aggregate_id = $1
        ORDER BY version DESC
        LIMIT 1
    `, aggregateID).Scan(
        &snapshot.SnapshotID,
        &snapshot.AggregateID,
        &snapshot.AggregateType,
        &snapshot.Version,
        &snapshot.State,
        &snapshot.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &snapshot, nil
} 
