{{ config(materialized='table') }}

with source as (

    select *
    from {{ source('raw', 'promotions') }}

),

cleaned as (

    select
        tenant_id,
        promo_id,
        promo_name,
        promo_type,
        discount_pct,
        start_date,
        end_date,
        min_purchase,
        created_at
    from source

)

select *
from cleaned