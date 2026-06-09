create database if not exists minimarket;

create table if not exists minimarket.raw_customers (
    customer_id Int32,
    name String,
    phone Nullable(String),
    email Nullable(String),
    gender Nullable(String),
    city Nullable(String),
    created_at datetime,
    loaded_at datetime
)
engine = MergeTree
order by customer_id;

create table if not exists minimarket.raw_products (
    product_id Int32,
    product_name String,
    category Nullable(String),
    brand Nullable(String),
    unit_price decimal(12, 2),
    is_active bool,
    created_at datetime,
    loaded_at datetime
)
engine = MergeTree
order by product_id;

create table if not exists minimarket.raw_transactions (
    transaction_id Int32,
    customer_id Nullable(Int32),
    store_id Nullable(Int32),
    transaction_date datetime,
    total_amount decimal(14, 2),
    payment_method Nullable(String),
    status Nullable(String),
    loaded_at datetime
)
engine = MergeTree
order by transaction_id;

create table if not exists minimarket.raw_transaction_items (
    item_id Int32,
    transaction_id Nullable(Int32),
    product_id Nullable(Int32),
    quantity Int32,
    unit_price decimal(12, 2),
    discount decimal(5, 2),
    subtotal decimal(14, 2),
    loaded_at datetime
)
engine = MergeTree
order by item_id;