package databases

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"minimarket-go-pipeline/internal/config"
	"time"

	_ "github.com/lib/pq"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

func NewPostgresConnection(cfg config.PostgresConfig) (*sql.DB, error) {
	fmt.Printf("===== DSN: %s\n", cfg.DSN())

	log.Default().Println("Connecting to Postgres Database")
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	log.Default().Println("Postgres connection successful")

	return db, nil
}

func NewClickHouseConnection(cfg config.ClickHouseConfig) (clickhouse.Conn, error) {
	fmt.Printf("===== Host:%s, Port:%s,Database:%s, User:%s, Password:%s\n", cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password)
	log.Default().Println("Connecting to ClickHouse Database")
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialTimeout:     10 * time.Second,
		ConnMaxLifetime: time.Hour,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	log.Default().Println("ClickHouse connection successful")

	return conn, nil
}
