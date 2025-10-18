DROP INDEX IF EXISTS idx_notifications_comment_id;
DROP INDEX IF EXISTS idx_notifications_post_id;
DROP INDEX IF EXISTS idx_notifications_user_id_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id_created_at;

DROP TABLE IF EXISTS notifications;
