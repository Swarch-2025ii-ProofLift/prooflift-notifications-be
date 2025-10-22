package repositories

import (
	"context"

	"prooflift-notifications-be/internal/models"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Save(ctx context.Context, input models.CreateNotificationInput) (*models.Notification, error)
	GetByID(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsReadByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteCommentNotification(ctx context.Context, commentID uuid.UUID) (int64, error)
	DeleteReactionNotification(ctx context.Context, postID uuid.UUID, actorID uuid.UUID) (int64, error)
}
