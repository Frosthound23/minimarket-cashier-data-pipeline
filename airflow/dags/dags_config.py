import os


def require_env(name: str) -> str:

    value = os.getenv(name)

    if not value:
        raise ValueError(f"Missing required environment variable: {name}")

    return value


PROJECT_ROOT = require_env("PROJECT_ROOT")
NETWORK_NAME = os.getenv("MINIMARKET_NETWORK_NAME", "minimarket_network")


POSTGRES_ENV = {
    "POSTGRES_HOST": require_env("POSTGRES_HOST"),
    "POSTGRES_PORT": os.getenv("POSTGRES_PORT", "5432"),
    "POSTGRES_DB": require_env("POSTGRES_DB"),
    "POSTGRES_USER": require_env("POSTGRES_USER"),
    "POSTGRES_PASSWORD": require_env("POSTGRES_PASSWORD"),
}


CLICKHOUSE_ENV = {
    "CLICKHOUSE_HOST": require_env("CLICKHOUSE_HOST"),
    "CLICKHOUSE_PORT": os.getenv("CLICKHOUSE_PORT", "8123"),
    "CLICKHOUSE_DB": require_env("CLICKHOUSE_DB"),
    "CLICKHOUSE_USER": require_env("CLICKHOUSE_USER"),
    "CLICKHOUSE_PASSWORD": require_env("CLICKHOUSE_PASSWORD"),
}
NETWORK_NAME = os.getenv("MINIMARKET_NETWORK_NAME", "minimarket_network")
PROJECT_ROOT = os.getenv("PROJECT_ROOT")