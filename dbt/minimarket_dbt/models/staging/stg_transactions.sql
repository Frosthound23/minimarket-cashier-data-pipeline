with source as (

    select *
    from {{ source('raw', 'raw_transactions') }}

),

cleaned as (

    select
        transaction_id,
        customer_id,
        store_id,
        transaction_date,
        total_amount,
        lower(payment_method) as payment_method,
        lower(status) as status,
        loaded_at
    from source
    where lower(status) = 'completed'

)

select *
from cleaned