package main

import (
	"log"

	"github.com/joho/godotenv"

	"minimarket-pipeline/internal/config"
	"minimarket-pipeline/internal/databases"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadAppConfig()

	postgresDB, err := databases.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	log.Println("PostgreSQL connection successful")

	clickhouseConn, err := databases.NewClickHouseConnection(cfg.ClickHouse)
	if err != nil {
		log.Fatal(err)
	}
	defer clickhouseConn.Close()

	log.Println("ClickHouse connection successful")

	tenants, err := config.LoadTenants("config/tenants.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, tenant := range tenants {
		log.Printf(
			"loaded tenant: tenant_id=%s schema=%s store_name=%s",
			tenant.TenantID,
			tenant.Schema,
			tenant.StoreName,
		)
	}
}
