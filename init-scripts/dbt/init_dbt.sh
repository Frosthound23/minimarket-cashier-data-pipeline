#!/bin/sh

set -e

echo "Installing dbt..."
pip install --no-cache-dir \
  "dbt-core==1.10.22" \
  "dbt-clickhouse==1.10.0"

echo "Current folder:"
pwd

echo "Files in /app:"
ls -la /app

echo "Files in /app/scripts:"
ls -la /app/scripts

echo "Running dbt debug..."
dbt debug --profiles-dir /app

echo "Running dbt run..."
dbt run --profiles-dir /app

echo "Running dbt test..."
dbt test --profiles-dir /app

echo "dbt completed successfully."