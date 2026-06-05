FROM apache/airflow:2.7.3-python3.10

ARG AIRFLOW_VERSION=2.7.3
ARG PYTHON_VERSION=3.10

USER airflow

RUN pip install --no-cache-dir \
    "apache-airflow-providers-docker" \
    "apache-airflow-providers-common-sql" \
    --constraint "https://raw.githubusercontent.com/apache/airflow/constraints-${AIRFLOW_VERSION}/constraints-${PYTHON_VERSION}.txt"
