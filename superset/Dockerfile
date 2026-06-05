FROM apache/superset:latest

USER root

RUN pip install --no-cache-dir \
    "clickhouse-connect>=0.8.0"

USER superset
