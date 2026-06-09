#!/bin/sh
set -e

echo "Initializing Superset..."

superset db upgrade

superset fab create-admin \
    --username admin \
    --firstname Admin \
    --lastname User \
    --email admin@example.com \
    --password admin



superset init

echo "Superset initialization complete!"