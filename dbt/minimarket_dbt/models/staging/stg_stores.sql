{{ config(materialized='table') }}

with source as (

    select *
    from {{ source('raw', 'stores') }}

),

cleaned as (

    select
        tenant_id,
        store_id,
        store_name,
        city,
        province,
        store_type,
        opened_at,
        is_active,
        loaded_at
    from source

)

select *
from cleaned