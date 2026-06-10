{{ config(materialized='table') }}

with source as (

    select *
    from {{ source('raw', 'products') }}

),

cleaned as (

    select
        tenant_id,
        product_id,
        product_name,
        category,
        brand,
        unit_price,
        is_active,
        created_at,
        loaded_at
    from source

)

select *
from cleaned