with source as (

    select *
    from {{ source('raw', 'raw_customers') }}

),

cleaned as (

    select
        customer_id,
        trim(name) as customer_name,
        phone,
        lower(email) as email,
        lower(gender) as gender,
        city,
        created_at,
        loaded_at
    from source

)

select *
from cleaned