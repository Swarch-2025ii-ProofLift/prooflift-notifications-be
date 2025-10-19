ALTER TABLE notifications DROP CONSTRAINT IF EXISTS chk_notification_type;

ALTER TABLE notifications ADD CONSTRAINT chk_notification_type
    CHECK (type IN ('COMMENT_CREATED', 'COMMENT_DELETED', 'REACTION_ADDED', 'REACTION_REMOVED'));

CREATE INDEX idx_notifications_post_actor_type ON notifications(post_id, actor_id, type);
