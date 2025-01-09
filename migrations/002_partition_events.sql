-- First, modify the events table to support partitioning
CREATE TABLE events (
    event_id        UUID NOT NULL,
    aggregate_id    UUID NOT NULL,
    aggregate_type  VARCHAR(255) NOT NULL,
    event_type      VARCHAR(255) NOT NULL,
    version         INTEGER NOT NULL,
    data            JSONB NOT NULL,
    metadata        JSONB,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Create partitions by month
CREATE TABLE events_y2024m01 PARTITION OF events
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE events_y2024m02 PARTITION OF events
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- Add constraints to partitions
ALTER TABLE events_y2024m01 ADD PRIMARY KEY (event_id);
ALTER TABLE events_y2024m01 ADD CONSTRAINT events_y2024m01_version_unique 
    UNIQUE (aggregate_id, version);

-- Create indexes on the parent table (automatically inherited by partitions)
CREATE INDEX idx_events_aggregate_lookup ON events(aggregate_id, version, created_at);
CREATE INDEX idx_events_type_lookup ON events(aggregate_type, created_at); 
