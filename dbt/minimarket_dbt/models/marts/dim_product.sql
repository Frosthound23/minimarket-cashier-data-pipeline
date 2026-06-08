select
    product_id as product_key,
    product_id,
    product_name,
    category,
    brand,
    unit_price,
    is_active,
    created_at,
    loaded_at
from {{ ref('stg_products') }}