from datetime import datetime, timezone
from typing import Any

import polars as pl
from sqlalchemy import text

from clients.databases import get_clickhouse_client, get_postgres_engine
from settings import settings
from loggings import logger

TABLE_MAPPINGS = {
    "customers": "raw_customers",
    "products": "raw_products",
    "transactions": "raw_transactions",
    "transaction_items": "raw_transaction_items",
}

def extract_table(postgres_engine, source_table: str) -> pl.DataFrame:
    query = text(f"select * from {source_table}")
    
    logger.info("Extracting data from %s", source_table)
    
    
    with postgres_engine.connect() as conn:
        result = conn.execute(query)
        rows = result.fetchall()
        columns = list(result.keys())
    
    df = pl.DataFrame(rows, schema=columns, orient="row")
    return df

def load_table_clickhouse(clickhouse_client: Any, df: pl.DataFrame, target_table: str) -> None:
    logger.info("Loading data into %s", target_table)

    if df.is_empty():
        logger.info("No data to load into %s. Skipping process.", target_table)
        return
    
    loaded_at = datetime.now(timezone.utc)
    
    df = df.with_columns(
        pl.lit(loaded_at).alias("loaded_at")
    )
    
    rows = df.rows()
    column_names = df.columns
    
    logger.info("Truncating Clickhouse Table: %s", target_table)
    clickhouse_client.command(f"truncate table {target_table}")
    
    logger.info("Loading %d rows into %s", df.height, target_table)
    clickhouse_client.insert(
        table=target_table, 
        data=rows, 
        column_names=column_names,
        )
    
    logger.info("Finished loading data %d into %s", df.height, target_table)
    
def run_pipeline():
    postgres_engine = get_postgres_engine()
    clickhouse_client = get_clickhouse_client()
    
    for source_table, target_table in TABLE_MAPPINGS.items():
        df = extract_table(postgres_engine, source_table)
        load_table_clickhouse(clickhouse_client, df, target_table)
        
    logger.info("Pipeline execution completed.")
    
    
if __name__ == "__main__":
    run_pipeline()
    
