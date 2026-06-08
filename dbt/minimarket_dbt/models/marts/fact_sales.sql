with transactions as (

    select *
    from {{ ref('stg_transactions') }}

),

transaction_items as (

    select *
    from {{ ref('stg_transaction_items') }}

),

final as (

    select
        ti.item_id as sales_id,
        t.transaction_id,
        ti.item_id,
        t.customer_id as customer_key,
        ti.product_id as product_key,
        toYYYYMMDD(toDate(t.transaction_date)) as date_key,
        t.store_id,
        t.transaction_date,
        t.payment_method,
        ti.quantity,
        ti.unit_price,
        ti.discount,
        ti.subtotal,
        t.loaded_at
    from transaction_items ti
    inner join transactions t
        on ti.transaction_id = t.transaction_id

)

select *
from final