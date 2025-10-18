package handlers

import (
	"net/http"
	"strconv"

	"prooflift-notifications-be/internal/api/middleware"
	"prooflift-notifications-be/internal/api/responses"
	"prooflift-notifications-be/internal/api/schemas"
	"prooflift-notifications-be/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		responses.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	notifications, err := h.service.ListUserNotifications(r.Context(), userID, limit, offset)
	if err != nil {
		responses.RespondError(w, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}

	response := schemas.NotificationList{
		Notifications: make([]schemas.Notification, 0, len(notifications)),
		Limit:         limit,
		Offset:        offset,
	}

	for _, n := range notifications {
		response.Notifications = append(response.Notifications, schemas.ToNotification(n))
	}

	responses.RespondSuccess(w, http.StatusOK, response)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		responses.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notificationIDStr := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		responses.RespondError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	notification, err := h.service.GetNotificationByID(r.Context(), notificationID)
	if err != nil {
		responses.RespondError(w, http.StatusNotFound, "notification not found")
		return
	}

	if notification.UserID != userID {
		responses.RespondError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.service.MarkAsRead(r.Context(), notificationID); err != nil {
		responses.RespondError(w, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	responses.RespondSuccess(w, http.StatusOK, map[string]string{"message": "notification marked as read"})
}

func (h *NotificationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.RespondSuccess(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "prooflift-notifications",
	})
}
