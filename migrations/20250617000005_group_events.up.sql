ALTER TABLE events ADD COLUMN group_id INT REFERENCES groups(id) ON DELETE CASCADE;

CREATE INDEX idx_events_group_start ON events(group_id, start_time) WHERE group_id IS NOT NULL;
