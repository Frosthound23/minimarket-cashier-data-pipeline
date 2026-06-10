{{ config(materialized='table') }}

select
    concat(tenant_id, '-', toString(customer_id)) as customer_key,
    tenant_id,
    customer_id,
    customer_name,
    phone,
    email,
    gender,
    city,
    created_at
from {{ ref('stg_customers') }}