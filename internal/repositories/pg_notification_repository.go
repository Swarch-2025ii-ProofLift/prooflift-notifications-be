package repositories

import (
	"context"
	"fmt"
	"time"

	"prooflift-notifications-be/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGNotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) NotificationRepository {
	return &PGNotificationRepository{db: db}
}

func (r *PGNotificationRepository) Save(ctx context.Context, input models.CreateNotificationInput) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, actor_id, post_id, comment_id, type, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, actor_id, post_id, comment_id, type, message, is_read, created_at, read_at
	`

	notification := &models.Notification{
		CreatedAt: time.Now(),
		IsRead:    false,
	}

	err := r.db.QueryRow(
		ctx,
		query,
		input.UserID,
		input.ActorID,
		input.PostID,
		input.CommentID,
		input.Type,
		input.Message,
		notification.IsRead,
		notification.CreatedAt,
	).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.ActorID,
		&notification.PostID,
		&notification.CommentID,
		&notification.Type,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
		&notification.ReadAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (r *PGNotificationRepository) GetByID(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, post_id, comment_id, type, message, is_read, created_at, read_at
		FROM notifications
		WHERE id = $1
	`

	notification := &models.Notification{}
	err := r.db.QueryRow(ctx, query, notificationID).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.ActorID,
		&notification.PostID,
		&notification.CommentID,
		&notification.Type,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
		&notification.ReadAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notification: %w", err)
	}

	return notification, nil
}

func (r *PGNotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, post_id, comment_id, type, message, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.ActorID,
			&notification.PostID,
			&notification.CommentID,
			&notification.Type,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *PGNotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, time.Now(), notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (r *PGNotificationRepository) MarkAllAsReadByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE user_id = $2 AND is_read = FALSE
	`
	result, err := r.db.Exec(ctx, query, time.Now(), userID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return result.RowsAffected(), nil
}

func (r *PGNotificationRepository) DeleteCommentNotification(ctx context.Context, commentID uuid.UUID) (int64, error) {
	query := `
		DELETE FROM notifications
		WHERE comment_id = $1
	`
	result, err := r.db.Exec(ctx, query, commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete notifications by comment_id: %w", err)
	}
	return result.RowsAffected(), nil
}

func (r *PGNotificationRepository) DeleteReactionNotification(ctx context.Context, postID uuid.UUID, actorID uuid.UUID) (int64, error) {
	query := `
		DELETE FROM notifications
		WHERE post_id = $1 AND actor_id = $2 AND type = 'REACTION_ADDED'
	`
	result, err := r.db.Exec(ctx, query, postID, actorID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete notifications by post_id and actor_id: %w", err)
	}
	return result.RowsAffected(), nil
}
