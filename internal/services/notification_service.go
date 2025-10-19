package services

import (
	"context"
	"fmt"

	"prooflift-notifications-be/internal/configs"
	"prooflift-notifications-be/internal/models"
	"prooflift-notifications-be/internal/repositories"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, actorID uuid.UUID, postID uuid.UUID, commentID *uuid.UUID, notfType models.NotificationType, message string) (*models.Notification, error) {
	if userID == uuid.Nil {
		return nil, configs.NewValidationError("userID cannot be nil")
	}
	if actorID == uuid.Nil {
		return nil, configs.NewValidationError("actorID cannot be nil")
	}
	if postID == uuid.Nil {
		return nil, configs.NewValidationError("postID cannot be nil")
	}

	if userID == actorID {
		return nil, configs.NewValidationError("cannot notify self")
	}

	switch notfType {
	case models.NotificationCommentCreated, models.NotificationReactionAdded:

	default:
		return nil, configs.NewValidationError(fmt.Sprintf("invalid notification type: %s", notfType))
	}

	input := models.CreateNotificationInput{
		UserID:    userID,
		ActorID:   actorID,
		PostID:    postID,
		CommentID: commentID,
		Type:      notfType,
		Message:   message,
	}

	notification, err := s.repo.Save(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	return notification, nil
}

func (s *NotificationService) ListUserNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error) {
	return s.repo.GetByID(ctx, notificationID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, notificationID)

	if err != nil {
		return fmt.Errorf("failed to fetch notification: %w", err)
	}

	if notification.IsRead {
		return nil
	}

	return s.repo.MarkAsRead(ctx, notificationID)
}

func (s *NotificationService) CreateCommentNotification(ctx context.Context, recipientID, actorID, postID uuid.UUID, commentID *uuid.UUID, message string) (*models.Notification, error) {
	return s.CreateNotification(ctx, recipientID, actorID, postID, commentID, models.NotificationCommentCreated, message)
}

func (s *NotificationService) CreateReactionNotification(ctx context.Context, recipientID, actorID, postID uuid.UUID, message string) (*models.Notification, error) {
	return s.CreateNotification(ctx, recipientID, actorID, postID, nil, models.NotificationReactionAdded, message)
}

func (s *NotificationService) DeleteCommentNotification(ctx context.Context, commentID uuid.UUID) (int64, error) {
	if commentID == uuid.Nil {
		return 0, configs.NewValidationError("commentID cannot be nil")
	}

	count, err := s.repo.DeleteCommentNotification(ctx, commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete notifications by comment_id: %w", err)
	}
	return count, nil
}

func (s *NotificationService) DeleteReactionNotifications(ctx context.Context, postID uuid.UUID, actorID uuid.UUID) (int64, error) {
	if postID == uuid.Nil {
		return 0, configs.NewValidationError("postID cannot be nil")
	}
	if actorID == uuid.Nil {
		return 0, configs.NewValidationError("actorID cannot be nil")
	}

	count, err := s.repo.DeleteReactionNotification(ctx, postID, actorID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete reaction notifications: %w", err)
	}
	return count, nil
}
