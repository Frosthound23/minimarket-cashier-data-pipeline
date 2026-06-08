import os
from datetime import datetime, timedelta

from airflow import DAG
from airflow.providers.docker.operators.docker import DockerOperator
from docker.types import Mount


NETWORK_NAME = os.getenv("MINIMARKET_NETWORK_NAME", "minimarket_network")
PROJECT_ROOT = os.getenv("PROJECT_ROOT")
if not PROJECT_ROOT:
    raise ValueError("PROJECT_ROOT environment variable is not set")


default_args = {
    "owner": "farrel",
    "description": "A DAG to orchestrate data pipeline for minimarket cashier data",
    "retries": 1,
    "retry_delay": timedelta(minutes=1),
}

with DAG(
    dag_id="minimarket_elt_pipeline",
    default_args=default_args,
    start_date=datetime(2026, 1, 1),
    schedule="0 15 * * *",
    catchup=False,
    tags=["minimarket", "elt", "dbt", "clickhouse"],
) as dag:
    run_pipeline = DockerOperator(
        task_id="run_python_polars_el",
        image="minimarket_pipeline:latest",
        command="python main.py",
        docker_url="unix://var/run/docker.sock",
        network_mode=NETWORK_NAME,
        auto_remove=True,
        mount_tmp_dir=False,
        mounts=[
            Mount(
                source=f"{PROJECT_ROOT}/pipeline",
                target="/app",
                type='bind',
                read_only=False,
            ),
        ],
        working_dir="/app",
    )

    dbt_debug = DockerOperator(
        task_id="dbt_debug",
        image="minimarket_dbt:latest",
        command="dbt debug --profiles-dir /usr/app/minimarket_dbt",
        docker_url="unix://var/run/docker.sock",
        network_mode=NETWORK_NAME,
        auto_remove=True,
        mount_tmp_dir=False,
        # environment=CLICKHOUSE_ENV,
        mounts=[
            Mount(
                source=f"{PROJECT_ROOT}/dbt/minimarket_dbt",
                target="/usr/app/minimarket_dbt",
                type="bind",
                read_only=False,
            ),
        ],
        working_dir="/usr/app/minimarket_dbt",
    )

    dbt_run = DockerOperator(
        task_id="dbt_run",
        image="minimarket_dbt:latest",
        command="dbt run --profiles-dir /usr/app/minimarket_dbt",
        docker_url="unix://var/run/docker.sock",
        network_mode=NETWORK_NAME,
        auto_remove=True,
        mount_tmp_dir=False,
        # environment=CLICKHOUSE_ENV,
        mounts=[
            Mount(
                source=f"{PROJECT_ROOT}/dbt/minimarket_dbt",
                target="/usr/app/minimarket_dbt",
                type="bind",
                read_only=False,
            ),
        ],
        working_dir="/usr/app/minimarket_dbt",
    )

    dbt_test = DockerOperator(
        task_id="dbt_test",
        image="minimarket_dbt:latest",
        command="dbt test --profiles-dir /usr/app/minimarket_dbt",
        docker_url="unix://var/run/docker.sock",
        network_mode=NETWORK_NAME,
        auto_remove=True,
        mount_tmp_dir=False,
        # environment=CLICKHOUSE_ENV,
        mounts=[
            Mount(
                source=f"{PROJECT_ROOT}/dbt/minimarket_dbt",
                target="/usr/app/minimarket_dbt",
                type="bind",
                read_only=False,
            ),
        ],
        working_dir="/usr/app/minimarket_dbt",
    )

    run_pipeline >> dbt_debug >> dbt_run >> dbt_test