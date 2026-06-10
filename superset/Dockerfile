FROM apache/superset:latest

USER root

RUN pip install --no-cache-dir "clickhouse-connect>=0.13.0"

USER superset
ENV SUPERSET_CONFIG_PATH=/app/pythonpath/superset_config.py