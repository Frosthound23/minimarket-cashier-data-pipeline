FROM python:3.10-slim

WORKDIR /usr/app

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DBT_PROFILES_DIR=/usr/app

RUN apt-get update \
    && apt-get install -y --no-install-recommends git build-essential \
    && pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir "dbt-core>=1.8,<1.9" "dbt-clickhouse>=1.8,<1.9" \
    && apt-get purge -y --auto-remove build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY . .

CMD ["dbt", "--version"]
