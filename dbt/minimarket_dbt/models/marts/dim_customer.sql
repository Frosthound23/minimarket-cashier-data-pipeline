select
    customer_id as customer_key,
    customer_id,
    customer_name,
    phone,
    email,
    gender,
    city,
    created_at,
    loaded_at
from {{ ref('stg_customers') }}