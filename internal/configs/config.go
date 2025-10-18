package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Service ServiceConfig
	DB      DatabaseConfig
	MQ      MessageQueueConfig
	HTTP    HTTPConfig
	JWT     JWTConfig
}

type ServiceConfig struct {
	Name      string
	QueueName string
}

type HTTPConfig struct {
	Port int
	Host string
}

type JWTConfig struct {
	Secret string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     int
	SSLMode  string
	URL      string
}

type MessageQueueConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	URL      string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Service: ServiceConfig{
			Name:      getEnv("SERVICE_NAME", "prooflift-notifications-be"),
			QueueName: getEnv("QUEUE_NAME", "notifications_queue"),
		},
		DB: DatabaseConfig{
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "root"),
			Host:     getEnv("DB_HOST", "prooflift-notifications-db"),
			Name:     getEnv("DB_NAME", "notifications_db"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		MQ: MessageQueueConfig{
			User:     getEnv("MQ_USER", "guest"),
			Password: getEnv("MQ_PASSWORD", "guest"),
			Host:     getEnv("MQ_HOST", "prooflift-notifications-mq"),
			Port:     getEnvAsInt("MQ_PORT", 5672),
		},
		HTTP: HTTPConfig{
			Port: getEnvAsInt("HTTP_PORT", 8080),
			Host: getEnv("HTTP_HOST", "0.0.0.0"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}

	cfg.DB.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode,
	)
	cfg.MQ.URL = fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.MQ.User, cfg.MQ.Password, cfg.MQ.Host, cfg.MQ.Port,
	)

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
