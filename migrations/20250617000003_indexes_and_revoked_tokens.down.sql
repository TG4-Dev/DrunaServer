DROP INDEX IF EXISTS idx_group_members_user;
DROP INDEX IF EXISTS idx_friends_friend_status;
DROP INDEX IF EXISTS idx_friends_user_status;
DROP INDEX IF EXISTS idx_events_user_start;
DROP INDEX IF EXISTS idx_revoked_tokens_expires_at;
DROP TABLE IF EXISTS revoked_tokens;
