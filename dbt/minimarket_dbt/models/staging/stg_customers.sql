{{ config(materialized='table') }}

with source as (

    select *
    from {{ source('raw', 'customers') }}

),

cleaned as (

    select
        tenant_id,
        customer_id,
        name as customer_name,
        phone,
        email,
        gender,
        city,
        created_at,
        loaded_at
    from source

)

select *
from cleaned