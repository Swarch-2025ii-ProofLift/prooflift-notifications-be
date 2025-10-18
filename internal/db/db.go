package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGConnection(url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour

	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := range maxRetries {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		db, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			cancel()
			if i < maxRetries-1 {
				log.Printf("Failed to create database connection pool (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("unable to create connection pool after %d attempts: %w", maxRetries, err)
		}

		if err = db.Ping(ctx); err != nil {
			db.Close()
			cancel()
			if i < maxRetries-1 {
				log.Printf("Failed to ping database (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("unable to ping database after %d attempts: %w", maxRetries, err)
		}

		cancel()
		log.Printf("Connected to database")
		return db, nil
	}

	return nil, fmt.Errorf("unable to connect to database after %d attempts", maxRetries)
}
