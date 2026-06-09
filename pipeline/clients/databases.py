import os
import sys
import time
from typing import Any
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
sys.path.append(str(PROJECT_ROOT))

from loggings import logger

import clickhouse_connect
from clickhouse_connect.driver.exceptions import OperationalError
from dotenv import load_dotenv
from sqlalchemy import create_engine
from sqlalchemy.engine import Engine
from settings import settings

load_dotenv()

def get_postgres_engine() -> Engine:
    postgres_url = settings.postgres_database_url
    
    return create_engine(postgres_url, pool_pre_ping=True)

def get_clickhouse_client(
    max_retries: int = 10,
    delay_seconds: int = 3,
) -> Any:
    logger.info("Starting connection to ClickHouse DB with max_retries=%d and delay_seconds=%d", max_retries, delay_seconds)
    last_error = None

    for attempt in range(1, max_retries + 1):
        try:
            client = clickhouse_connect.get_client(
                host=os.getenv("CLICKHOUSE_HOST"),
                port=int(os.getenv("CLICKHOUSE_PORT", "8123")),
                username=os.getenv("CLICKHOUSE_USER"),
                password=os.getenv("CLICKHOUSE_PASSWORD"),
                database=os.getenv("CLICKHOUSE_DB"),
            )

            client.command("SELECT 1")
            logger.info("Connected to ClickHouse successfully")
            return client

        except Exception as error:
            last_error = error
            print(
                f"ClickHouse not ready. "
                f"Attempt {attempt}/{max_retries}. "
                f"Retrying in {delay_seconds} seconds..."
            )
            time.sleep(delay_seconds)

    raise ConnectionError(
        "Failed to connect to ClickHouse after retries"
    ) from last_error
    
get_postgres_engine()
get_clickhouse_client()