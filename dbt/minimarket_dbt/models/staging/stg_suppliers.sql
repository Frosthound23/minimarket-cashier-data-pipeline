{{ config(materialized='table') }}

with source as (

    select *
    from {{ source('raw', 'suppliers') }}

),

cleaned as (

    select
        tenant_id,
        supplier_id,
        supplier_name,
        contact_name,
        city,
        country,
        created_at,
        loaded_at
    from source

)

select *
from cleaned