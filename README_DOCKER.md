# Docker setup

## Start the infrastructure

```bash
docker compose up -d --build
```

## Manual checks

```bash
docker compose ps

docker compose run --rm dbt dbt debug
docker compose run --rm dbt dbt run
docker compose run --rm dbt dbt test

docker compose run --rm elt
```

## Services

- Airflow: http://localhost:8080
  - username: admin
  - password: admin
- Superset: http://localhost:8088
  - username: admin
  - password: admin
- ClickHouse HTTP: http://localhost:8123
- PostgreSQL source: localhost:5432

## Superset ClickHouse connection URI

Use this URI in Superset database connection:

```text
clickhousedb://default:@clickhouse:8123/pos_dwh
```

If your local Superset driver asks for SQLAlchemy URI and the previous one fails, try:

```text
clickhousedb+connect://default:@clickhouse:8123/pos_dwh
```

## Apple Silicon note

Do not add `platform: linux/amd64` by default. Add it only to the service that fails with architecture-related errors.
