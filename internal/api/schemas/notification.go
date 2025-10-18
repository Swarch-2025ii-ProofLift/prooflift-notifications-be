package schemas

import (
	"time"

	"prooflift-notifications-be/internal/models"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ActorID   uuid.UUID  `json:"actor_id"`
	PostID    uuid.UUID  `json:"post_id"`
	CommentID *uuid.UUID `json:"comment_id,omitempty"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

type NotificationList struct {
	Notifications []Notification `json:"notifications"`
	Limit         int            `json:"limit"`
	Offset        int            `json:"offset"`
}

func ToNotification(n *models.Notification) Notification {
	return Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		ActorID:   n.ActorID,
		PostID:    n.PostID,
		CommentID: n.CommentID,
		Type:      string(n.Type),
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		ReadAt:    n.ReadAt,
	}
}
