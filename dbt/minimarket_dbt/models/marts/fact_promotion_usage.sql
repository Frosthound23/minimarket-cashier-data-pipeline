{{ config(materialized='table') }}

select
    concat(tp.tenant_id, '-', toString(tp.transaction_promotion_id)) as promotion_usage_key,

    tp.tenant_id as tenant_id,

    concat(tp.tenant_id, '-', toString(tp.transaction_id)) as transaction_key,
    concat(tp.tenant_id, '-', toString(tp.promo_id)) as promotion_key,

    tp.transaction_promotion_id as transaction_promotion_id,
    tp.transaction_id as transaction_id,
    tp.promo_id as promo_id,

    t.customer_id as customer_id,
    t.store_id as store_id,

    concat(t.tenant_id, '-', toString(t.customer_id)) as customer_key,
    concat(t.tenant_id, '-', toString(t.store_id)) as store_key,

    t.transaction_date as transaction_date,
    t.transaction_date_key as transaction_date_key,

    p.promo_name as promo_name,
    p.promo_type as promo_type,
    p.discount_pct as discount_pct,
    p.min_purchase as min_purchase,

    tp.discount_applied as discount_applied,
    t.total_amount as total_amount,

    t.loaded_at as transaction_loaded_at,
    tp.loaded_at as promotion_usage_loaded_at
from {{ ref('stg_transaction_promotions') }} tp
inner join {{ ref('stg_transactions') }} t
    on tp.tenant_id = t.tenant_id
    and tp.transaction_id = t.transaction_id
inner join {{ ref('stg_promotions') }} p
    on tp.tenant_id = p.tenant_id
    and tp.promo_id = p.promo_id