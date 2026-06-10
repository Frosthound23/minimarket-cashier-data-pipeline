{{ config(materialized='table') }}

select
    concat(ti.tenant_id, '-', toString(ti.item_id)) as sales_key,

    ti.tenant_id,

    concat(t.tenant_id, '-', toString(t.transaction_id)) as transaction_key,
    concat(t.tenant_id, '-', toString(t.customer_id)) as customer_key,
    concat(t.tenant_id, '-', toString(t.store_id)) as store_key,
    concat(ti.tenant_id, '-', toString(ti.product_id)) as product_key,

    t.transaction_id,
    ti.item_id,
    t.customer_id,
    t.store_id,
    ti.product_id,

    t.transaction_date,
    t.transaction_date_key,

    t.payment_method,
    t.status,

    ti.quantity,
    ti.unit_price,
    ti.discount,
    ti.subtotal,

    t.total_amount
from {{ ref('stg_transaction_items') }} ti
inner join {{ ref('stg_transactions') }} t
    on ti.tenant_id = t.tenant_id
    and ti.transaction_id = t.transaction_id