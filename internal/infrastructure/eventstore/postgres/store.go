package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nedson202/go-cqrs/internal/domain/events"
)

type PostgresEventStore struct {
	db        *sql.DB
	batchSize int
}

func NewEventStore(db *sql.DB, batchSize int) *PostgresEventStore {
	return &PostgresEventStore{
		db:        db,
		batchSize: batchSize,
	}
}

// Save implements EventStore.Save
func (s *PostgresEventStore) Save(ctx context.Context, events []events.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("events",
		"event_id", "aggregate_id", "aggregate_type", "event_type",
		"version", "data", "metadata", "created_at"))
	if err != nil {
		return fmt.Errorf("prepare copy statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.ExecContext(ctx,
			event.EventID,
			event.AggregateID,
			event.AggregateType,
			event.EventType,
			event.Version,
			event.Data,
			event.Metadata,
			event.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("execute copy: %w", err)
		}
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("close copy statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Load implements EventStore.Load
func (s *PostgresEventStore) Load(ctx context.Context, aggregateID string, fromVersion int) ([]events.Event, error) {
	aggregateIDUUID, err := uuid.Parse(aggregateID)
	if err != nil {
		return nil, fmt.Errorf("parse aggregate ID: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT event_id, aggregate_id, aggregate_type, event_type,
			   version, data, metadata, created_at
		FROM events
		WHERE aggregate_id = $1 AND version > $2
		ORDER BY version ASC
	`, aggregateIDUUID, fromVersion)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var eventList []events.Event
	for rows.Next() {
		var event events.Event
		err := rows.Scan(
			&event.EventID,
			&event.AggregateID,
			&event.AggregateType,
			&event.EventType,
			&event.Version,
			&event.Data,
			&event.Metadata,
			&event.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		eventList = append(eventList, event)
	}

	return eventList, nil
}

// StreamEvents implements EventStore.StreamEvents
func (s *PostgresEventStore) Stream(ctx context.Context, aggregateID string, fromVersion int) (<-chan events.Event, <-chan error) {
	eventCh := make(chan events.Event)
	errCh := make(chan error, 1)

	aggregateIDUUID, err := uuid.Parse(aggregateID)
	if err != nil {
		errCh <- fmt.Errorf("parse aggregate ID: %w", err)
		return eventCh, errCh
	}

	go func() {
		defer close(eventCh)
		defer close(errCh)

		rows, err := s.db.QueryContext(ctx, `
			SELECT event_id, aggregate_id, aggregate_type, event_type,
				   version, data, metadata, created_at
			FROM events
			WHERE aggregate_id = $1 AND version > $2
			ORDER BY version ASC
		`, aggregateIDUUID, fromVersion)
		if err != nil {
			errCh <- fmt.Errorf("query events: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var event events.Event
			err := rows.Scan(
				&event.EventID,
				&event.AggregateID,
				&event.AggregateType,
				&event.EventType,
				&event.Version,
				&event.Data,
				&event.Metadata,
				&event.Timestamp,
			)
			if err != nil {
				errCh <- fmt.Errorf("scan event: %w", err)
				return
			}

			select {
			case eventCh <- event:
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventCh, errCh
}

// Ensure PostgresEventStore implements EventStore
var _ events.EventStore = (*PostgresEventStore)(nil) 
