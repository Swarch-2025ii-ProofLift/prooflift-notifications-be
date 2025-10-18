package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationCommentCreated NotificationType = "COMMENT_CREATED"
	NotificationReactionAdded  NotificationType = "REACTION_ADDED"
)

type Notification struct {
	ID        uuid.UUID        `db:"id"`
	UserID    uuid.UUID        `db:"user_id"`
	ActorID   uuid.UUID        `db:"actor_id"`
	PostID    uuid.UUID        `db:"post_id"`
	Type      NotificationType `db:"type"`
	CommentID *uuid.UUID       `db:"comment_id"`
	Message   string           `db:"message"`
	IsRead    bool             `db:"is_read"`
	CreatedAt time.Time        `db:"created_at"`
	ReadAt    *time.Time       `db:"read_at"`
}

type CreateNotificationInput struct {
	UserID    uuid.UUID
	ActorID   uuid.UUID
	PostID    uuid.UUID
	CommentID *uuid.UUID
	Type      NotificationType
	Message   string
}
