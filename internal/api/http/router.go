package http

import (
	"time"

	"prooflift-notifications-be/internal/api/handlers"
	"prooflift-notifications-be/internal/api/middleware"
	"prooflift-notifications-be/internal/services"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(notificationService *services.NotificationService, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	notificationHandler := handlers.NewNotificationHandler(notificationService)

	r.Get("/health", notificationHandler.HealthCheck)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(jwtSecret))

		r.Get("/notifications", notificationHandler.ListNotifications)
		r.Put("/notifications/read-all", notificationHandler.MarkAllAsRead)
		r.Put("/notifications/{id}/read", notificationHandler.MarkAsRead)
	})

	return r
}
