#!/bin/sh
set -e
sleep 5
airflow db init
airflow users create \
    --username airflow \
    --firstname Admin \
    --lastname User \
    --role Admin \
    --email admin@example.com \
    --password airflow

echo "Airflow initialization complete!"