DROP INDEX IF EXISTS idx_notifications_post_actor_type;

ALTER TABLE notifications DROP CONSTRAINT IF EXISTS chk_notification_type;

ALTER TABLE notifications ADD CONSTRAINT chk_notification_type
    CHECK (type IN ('COMMENT_CREATED', 'REACTION_ADDED'));
