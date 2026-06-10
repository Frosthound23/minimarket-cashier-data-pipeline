with source as (

    select *
    from {{ source('raw', 'transaction_items') }}

),

cleaned as (

    select
        tenant_id,
        item_id,
        transaction_id,
        product_id,
        quantity,
        unit_price,
        discount,
        subtotal,
        loaded_at
    from source

)

select *
from cleaned