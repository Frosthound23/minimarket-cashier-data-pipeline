with source as (

    select *
    from {{ source('raw', 'raw_products') }}

),

cleaned as (

    select
        product_id,
        trim(product_name) as product_name,
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