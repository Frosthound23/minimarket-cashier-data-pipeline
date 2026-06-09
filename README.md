# Minimarket Cashier Data Pipeline Beginner - Level

## 1. Project Description

This project is an end-to-end data engineering pipeline for a minimarket cashier / point-of-sale system at a beginner-level. This project is for the take home test of the data engineering job application in Parkee.

The pipeline extracts transactional data from PostgreSQL, loads the raw data into ClickHouse, transforms the data using dbt into a dimensional star schema, orchestrates the workflow using Apache Airflow, and visualizes the analytical results using Jupyter Notebook.

This project is built for the **Data Engineer Take-Home Technical Test** at the Beginner level.

## 2. Tech Stack

| Component         | Tool                      | Notes                                             |
| ----------------- | ------------------------- | ------------------------                          |
| Source Database   | PostgreSQL                | OLTP                                              |
| Data Warehouse    | ClickHouse                | OLAP                                              |
| Pipeline Language | Python                    |                                                   |
| DataFrame Library | Polars                    | Pandas Alternative because it is too slow         |
| Transformation    | dbt Core + dbt-clickhouse |                                       |  
| Orchestration     | Apache Airflow            |                                       |
| Visualization     | Jupyter Notebook          |                                                   |
| Optional BI Tool  | Apache Superset           | Not being used in beginner-level                  |
| Containerization  | Docker Compose            |                                                   |

## 3. Architecture

```mermaid
flowchart TD
    A[PostgreSQL Source Database] --> B[Python Polars EL Pipeline]
    B --> C[ClickHouse Raw Tables]
    C --> D[dbt Staging Models]
    D --> E[dbt Mart Models]
    E --> F[Jupyter Notebook Visualization]
    E --> G[Apache Superset (NOT BEING USED RIGHT NOW JUST FOR SAFE KEEPING in the Intermediate-Level)]

    H[Apache Airflow] --> B
    H --> D
    H --> E
```

## 4. Data Flow

```text
PostgreSQL
    ↓
Python Polars Extract & Load
    ↓
ClickHouse raw tables
    ↓
dbt staging models
    ↓
dbt mart models
    ↓
Jupyter Notebook analytics
```

## 5. Star Schema

The final mart layer follows a star schema.

```text
                    dim_date
                       |
dim_customer ---- fact_sales ---- dim_product
```

### Main mart tables

| Model          | Description                                     |
| -------------- | ----------------------------------------------- |
| `dim_customer` | Customer dimension to keep customer's data dimension         |
| `dim_product`  | Product dimension to keep every products                               |
| `dim_date`     | Date dimension generated from every transaction dates |
| `fact_sales`   | Sales fact table for every transaction items that was sold|

The grain of `fact_sales` is **one row per transaction item**. This allows product-level analysis, such as top products by category.

## 6. Analytical Questions
Can read in the analysis.ipynb in the folder minimarket-cashier-data-pipeline/notebooks/

## 7. Repository Structure

```text
minimarket-cashier-data-pipeline/
├── README.md
├── docker-compose.yml
├── .env
├── .venv(python's virtual environment but need set up)
├── .gitignore
├── requirements.txt
|
├── airflow/
│   ├── Dockerfile
|   ├── logs 
│   └── dags/
│       └── minimarket_elt_dags.py
|       └── dags_config.py
│
├── dbt/
│   ├── Dockerfile
│   ├── requirements.txt
|   ├── logs/
│   └── minimarket_dbt/
│       ├── dbt_project.yml
│       ├── profiles.yml
│       └── models/
│           ├── staging/
|           |   ├── sources.yml
|           |   ├── stg_customers.sql 
|           |   ├── stg_products.sql 
|           |   ├── stg_transaction_items.sql
|           |   └── stg_transactions.sql
│           └── marts/
|               ├── sources.yml
|               ├── dim_customer.sql
|               ├── dim_date.sql
|               ├── dim_product.sql
|               └── fact_sales.sql
|               
│
├── init-scripts/
|   ├── airflow/
|   |   └── init_airflow.sh
│   ├── postgres/
|   |   ├── seed_data.sql        
│   │   └── init_postgres.sql
│   ├── clickhouse/
│   │   └── init_clickhouse.sql
│   └── dbt/
|       └── init_dbt.sh
│
├── logs/
|
├── notebooks/
│   └── analysis.ipynb
│
├── pipeline/
│   ├── Dockerfile
│   ├── requirements.txt
|   ├── loggings.py
│   ├── main.py
│   ├── settings.py
|   ├── clients/
|   |   └── databases.py
|   └── config/
|       └── tenants.json
|   
│
└── superset/
    ├── Dockerfile
    └── requirements.txt
```

