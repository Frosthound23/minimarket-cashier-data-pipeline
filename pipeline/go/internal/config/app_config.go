package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Postgres   PostgresConfig
	ClickHouse ClickHouseConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
}

type ClickHouseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func LoadAppConfig() AppConfig {

	return AppConfig{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Database: getEnv("POSTGRES_DB", "minimarket"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		ClickHouse: ClickHouseConfig{
			Host:     getEnv("CLICKHOUSE_HOST", "clickhouse"),
			Port:     getEnv("CLICKHOUSE_NATIVE_PORT", "9000"),
			Database: getEnv("CLICKHOUSE_DB", "minimarket"),
			User:     getEnv("CLICKHOUSE_USER", "minimarket_user"),
			Password: getEnv("CLICKHOUSE_PASSWORD", "minimarket_password"),
		},
	}
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Database,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
