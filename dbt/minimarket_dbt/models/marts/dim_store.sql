{{ config(materialized='table') }}

select
    concat(tenant_id, '-', toString(store_id)) as store_key,
    tenant_id,
    store_id,
    store_name,
    city,
    province,
    store_type,
    opened_at,
    is_active
from {{ ref('stg_stores') }}