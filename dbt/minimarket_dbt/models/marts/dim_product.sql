{{ config(materialized='table') }}

select
    concat(tenant_id, '-', toString(product_id)) as product_key,
    tenant_id,
    product_id,
    product_name,
    category,
    brand,
    unit_price,
    is_active,
    created_at
from {{ ref('stg_products') }}