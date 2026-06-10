with source as (

    select *
    from {{ source('raw', 'transactions') }}

),

cleaned as (

    select
        tenant_id,
        transaction_id,
        customer_id,
        store_id,
        transaction_date,
        toDate(transaction_date) as transaction_date_key,
        total_amount,
        payment_method,
        status,
        created_at
    from source
    where lower(status) = 'completed'

)

select *
from cleaned