import os

import clickhouse_connect
from dotenv import load_dotenv

load_dotenv()


def get_clickhouse_client():
    return clickhouse_connect.get_client(
        host=os.getenv("CLICKHOUSE_HOST", "clickhouse"),
        port=int(os.getenv("CLICKHOUSE_PORT", "8123")),
        username=os.getenv("CLICKHOUSE_USER", "minimarket_user"),
        password=os.getenv("CLICKHOUSE_PASSWORD", "minimarket_password"),
        database=os.getenv("CLICKHOUSE_DB", "minimarket"),
    )