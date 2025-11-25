package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret     string
	GelfEndpoint  string
	ServerPort    string
	PostgresqlDSN string
	BrokerHost    string
}

const (
	JWT_SECRET        = "fallnacl-secret-change-in-production"
	GELF_ENDPOINT     = "graylog:12201"
	SERVER_PORT       = ":8080"
	POSTGRESQL_DSN    = "host=catalog-db-postgresql user=catalog_user password=catalog_pass dbname=catalog_db port=5432 sslmode=disable"
	KAFKA_BROKER_HOST = "kafka:9092"
)

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", JWT_SECRET),
		PostgresqlDSN: getEnv("DB_DSN", POSTGRESQL_DSN),
		GelfEndpoint:  getEnv("GELF_ENDPOINT", GELF_ENDPOINT),
		ServerPort:    getEnv("SERVER_PORT", SERVER_PORT),
		BrokerHost:    getEnv("KAFKA_BROKER_HOST", KAFKA_BROKER_HOST),
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
