#!/bin/sh
# Healthcheck
sleep 5

echo "Init ClickHouse databases and tables"

clickhouse-client --query="CREATE DATABASE IF NOT EXISTS dwh;"

clickhouse-client --query="
CREATE TABLE IF NOT EXISTS dwh.dim_customer (
    customer_id UInt32,
    name String,
    phone String,
    email String,
    gender String,
    city String,
    created_at DateTime
) ENGINE = MergeTree()
ORDER BY customer_id;
"

clickhouse-client --query="
CREATE TABLE IF NOT EXISTS dwh.dim_product (
    product_id UInt32,
    product_name String,
    category String,
    brand String,
    unit_price Float64,
    created_at DateTime
) ENGINE = MergeTree()
ORDER BY product_id;
"

clickhouse-client --query="
CREATE TABLE IF NOT EXISTS dwh.dim_date (
    date Date,
    year UInt16,
    month UInt8,
    day UInt8,
    weekday UInt8
) ENGINE = MergeTree()
ORDER BY date;
"

clickhouse-client --query="
CREATE TABLE IF NOT EXISTS dwh.fact_sales (
    transaction_id UInt32,
    customer_id UInt32,
    product_id UInt32,
    quantity UInt32,
    total_amount Float64,
    transaction_date DateTime
) ENGINE = MergeTree()
ORDER BY (transaction_id, customer_id);
"

echo "ClickHouse init success"