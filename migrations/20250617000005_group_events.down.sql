DROP INDEX IF EXISTS idx_events_group_start;

ALTER TABLE events DROP COLUMN IF EXISTS group_id;
