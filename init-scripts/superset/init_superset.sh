#!/bin/sh
set -e

echo "Initializing Superset..."

# Initialize database
superset db upgrade

# Create admin user
superset fab create-admin \
    --username admin \
    --firstname Admin \
    --lastname User \
    --email admin@example.com \
    --password admin

# Load examples (optional, remove if not needed)
# superset load_examples

# Initialize roles and permissions
superset init

echo "Superset initialization complete!"