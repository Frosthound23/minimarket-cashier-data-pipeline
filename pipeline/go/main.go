package main

import (
	"context"
	"log"
	"minimarket-go-pipeline/internal/app"
	"minimarket-go-pipeline/internal/config"
	"minimarket-go-pipeline/internal/databases"
	"minimarket-go-pipeline/internal/extractor"
	"minimarket-go-pipeline/internal/loader"
	"minimarket-go-pipeline/internal/watermark"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadAppConfig()

	tenants, err := config.LoadTenants("config/tenants.json")
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := databases.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	clickhouseConn, err := databases.NewClickHouseConnection(cfg.ClickHouse)
	if err != nil {
		log.Fatal(err)
	}
	defer clickhouseConn.Close()

	postgresExtractor := extractor.NewPostgresExtractor(postgresDB)
	clickhouseLoader := loader.NewClickHouseLoader(clickhouseConn)

	watermarkStore := watermark.NewStore(clickhouseConn)

	pipeline := app.NewPipeline(
		postgresExtractor,
		clickhouseLoader,
		watermarkStore,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := pipeline.Run(ctx, tenants); err != nil {
		log.Fatal(err)
	}

	log.Println("pipeline completed successfully")
}
