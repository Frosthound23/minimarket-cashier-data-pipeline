with source as (

    select *
    from {{ source('raw', 'transaction_promotions') }}

),

cleaned as (

    select
        tenant_id,
        id as transaction_promotion_id,
        transaction_id,
        promo_id,
        discount_applied,
        created_at
    from source

)

select *
from cleaned