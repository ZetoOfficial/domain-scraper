package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQQueue    string
}

func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Warning: could not load config from %s, using environment variables", path)
	}

	config := &Config{
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMQQueue:    getEnv("RABBITMQ_QUEUE", "links_queue"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
