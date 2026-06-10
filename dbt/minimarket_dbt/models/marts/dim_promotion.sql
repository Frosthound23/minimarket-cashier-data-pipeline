{{ config(materialized='table') }}

select
    concat(tenant_id, '-', toString(promo_id)) as promotion_key,
    tenant_id,
    promo_id,
    promo_name,
    promo_type,
    discount_pct,
    start_date,
    end_date,
    min_purchase,
    loaded_at
from {{ ref('stg_promotions') }}