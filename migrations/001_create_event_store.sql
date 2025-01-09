CREATE TABLE events (
    event_id        UUID PRIMARY KEY,
    aggregate_id    UUID NOT NULL,
    aggregate_type  VARCHAR(255) NOT NULL,
    event_type      VARCHAR(255) NOT NULL,
    version         INTEGER NOT NULL,
    data            JSONB NOT NULL,
    metadata        JSONB,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (aggregate_id, version)
);

CREATE TABLE snapshots (
    snapshot_id     UUID PRIMARY KEY,
    aggregate_id    UUID NOT NULL,
    aggregate_type  VARCHAR(255) NOT NULL,
    version         INTEGER NOT NULL,
    state          JSONB NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_events_aggregate ON events(aggregate_id, version);
CREATE INDEX idx_snapshots_aggregate ON snapshots(aggregate_id, version); 
