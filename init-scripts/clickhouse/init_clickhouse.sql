CREATE DATABASE IF NOT EXISTS minimarket;

CREATE TABLE IF NOT EXISTS minimarket.raw_customers (
    tenant_id String,
    customer_id Int32,
    name String,
    phone Nullable(String),
    email Nullable(String),
    gender Nullable(String),
    city Nullable(String),
    created_at DateTime,
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, customer_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_products (
    tenant_id String,
    product_id Int32,
    product_name String,
    category Nullable(String),
    brand Nullable(String),
    unit_price Decimal(12, 2),
    is_active Bool,
    created_at DateTime,
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, product_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_stores (
    tenant_id String,
    store_id Int32,
    store_name String,
    city Nullable(String),
    province Nullable(String),
    store_type Nullable(String),
    opened_at Nullable(Date),
    is_active Bool,
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, store_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_promotions (
    tenant_id String,
    promo_id Int32,
    promo_name Nullable(String),
    promo_type Nullable(String),
    discount_pct Nullable(Decimal(5, 2)),
    start_date Nullable(Date),
    end_date Nullable(Date),
    min_purchase Decimal(12, 2),
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, promo_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_suppliers (
    tenant_id String,
    supplier_id Int32,
    supplier_name Nullable(String),
    contact_name Nullable(String),
    city Nullable(String),
    country Nullable(String),
    created_at DateTime,
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, supplier_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_transactions (
    tenant_id String,
    transaction_id Int32,
    customer_id Nullable(Int32),
    store_id Nullable(Int32),
    transaction_date DateTime,
    total_amount Decimal(14, 2),
    payment_method Nullable(String),
    status Nullable(String),
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, transaction_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_transaction_items (
    tenant_id String,
    item_id Int32,
    transaction_id Nullable(Int32),
    product_id Nullable(Int32),
    quantity Int32,
    unit_price Decimal(12, 2),
    discount Decimal(12, 2),
    subtotal Decimal(14, 2),
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, item_id);

CREATE TABLE IF NOT EXISTS minimarket.raw_transaction_promotions (
    tenant_id String,
    id Int32,
    transaction_id Nullable(Int32),
    promo_id Nullable(Int32),
    discount_applied Nullable(Decimal(12, 2)),
    loaded_at DateTime
)
ENGINE = MergeTree
ORDER BY (tenant_id, id);

CREATE TABLE IF NOT EXISTS minimarket.elt_watermarks (
    tenant_id String,
    table_name String,
    watermark_column String,
    last_watermark DateTime,
    processed_at DateTime
)
ENGINE = ReplacingMergeTree(processed_at)
ORDER BY (tenant_id, table_name);