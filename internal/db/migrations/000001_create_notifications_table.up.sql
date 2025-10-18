CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    actor_id UUID NOT NULL,
    post_id UUID NOT NULL,
    comment_id UUID NULL,
    type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    read_at TIMESTAMP NULL
);

CREATE INDEX idx_notifications_user_id_created_at ON notifications(user_id, created_at DESC);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);

CREATE INDEX idx_notifications_post_id ON notifications(post_id);

CREATE INDEX idx_notifications_comment_id ON notifications(comment_id) WHERE comment_id IS NOT NULL;

ALTER TABLE notifications ADD CONSTRAINT chk_notification_type
    CHECK (type IN ('COMMENT_CREATED', 'REACTION_ADDED'));

ALTER TABLE notifications ADD CONSTRAINT chk_no_self_notification
    CHECK (user_id != actor_id);
