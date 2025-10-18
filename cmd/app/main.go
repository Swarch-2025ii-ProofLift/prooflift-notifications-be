package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customhttp "prooflift-notifications-be/internal/api/http"
	"prooflift-notifications-be/internal/configs"
	"prooflift-notifications-be/internal/db"
	"prooflift-notifications-be/internal/events"
	"prooflift-notifications-be/internal/mq"
	"prooflift-notifications-be/internal/repositories"
	"prooflift-notifications-be/internal/services"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("Running database migrations...")
	if err := db.RunMigrations(cfg.DB.URL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	dbPool, err := db.NewPGConnection(cfg.DB.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	mqConn, err := mq.NewMQConnection(cfg.MQ.URL)
	if err != nil {
		log.Fatalf("failed to connect to message queue: %v", err)
	}
	defer mqConn.Close()

	notificationRepo := repositories.NewNotificationRepository(dbPool)
	notificationService := services.NewNotificationService(notificationRepo)

	consumer, err := events.NewConsumer(mqConn, cfg.Service.QueueName, notificationService)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer consumer.Close()

	router := customhttp.NewRouter(notificationService, cfg.JWT.Secret)
	serverAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	httpServer := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Printf("Starting HTTP server on %s", serverAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting MQ consumer...")
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Received shutdown signal, stopping service...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	cancel()

	log.Println("Service stopped successfully")
}
