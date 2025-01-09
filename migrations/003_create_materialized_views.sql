-- Create materialized view for current aggregate state
CREATE MATERIALIZED VIEW aggregate_states AS
WITH latest_events AS (
    SELECT DISTINCT ON (aggregate_id)
        aggregate_id,
        aggregate_type,
        version,
        data,
        created_at
    FROM events
    ORDER BY aggregate_id, version DESC
)
SELECT 
    aggregate_id,
    aggregate_type,
    version,
    data as state,
    created_at as last_updated
FROM latest_events;

-- Create index for efficient lookups
CREATE UNIQUE INDEX idx_aggregate_states_id ON aggregate_states(aggregate_id);
CREATE INDEX idx_aggregate_states_type ON aggregate_states(aggregate_type); 
