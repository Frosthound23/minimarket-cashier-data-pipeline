{{ config(materialized='table') }}

select
    concat(tp.tenant_id, '-', toString(tp.transaction_promotion_id)) as promotion_usage_key,

    tp.tenant_id,

    concat(tp.tenant_id, '-', toString(tp.transaction_id)) as transaction_key,
    concat(tp.tenant_id, '-', toString(tp.promo_id)) as promotion_key,

    tp.transaction_promotion_id,
    tp.transaction_id,
    tp.promo_id,

    t.transaction_date,
    t.transaction_date_key,

    p.promo_name,
    p.promo_type,
    p.discount_pct,
    p.min_purchase,

    tp.discount_applied,
    t.total_amount
from {{ ref('stg_transaction_promotions') }} tp
inner join {{ ref('stg_transactions') }} t
    on tp.tenant_id = t.tenant_id
    and tp.transaction_id = t.transaction_id
inner join {{ ref('stg_promotions') }} p
    on tp.tenant_id = p.tenant_id
    and tp.promo_id = p.promo_id