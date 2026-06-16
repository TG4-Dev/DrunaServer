ALTER TABLE group_members
    ADD CONSTRAINT unique_group_member UNIQUE (group_id, user_id);
