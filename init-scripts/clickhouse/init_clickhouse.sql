create database if not exists minimarket;

create table if not exists minimarket.raw_customers (
    customer_id int32,
    name string,
    phone nullable(string),
    email nullable(string),
    gender nullable(string),
    city nullable(string),
    created_at datetime,
    loaded_at datetime
)
engine = MergeTree
order by customer_id;

create table if not exists minimarket.raw_products (
    product_id int32,
    product_name string,
    category nullable(string),
    brand nullable(string),
    unit_price decimal(12, 2),
    is_active bool,
    created_at datetime,
    loaded_at datetime
)
engine = MergeTree
order by product_id;

create table if not exists minimarket.raw_transactions (
    transaction_id int32,
    customer_id nullable(int32),
    store_id nullable(int32),
    transaction_date datetime,
    total_amount decimal(14, 2),
    payment_method nullable(string),
    status nullable(string),
    loaded_at datetime
)
engine = MergeTree
order by transaction_id;

create table if not exists minimarket.raw_transaction_items (
    item_id int32,
    transaction_id nullable(int32),
    product_id nullable(int32),
    quantity int32,
    unit_price decimal(12, 2),
    discount decimal(5, 2),
    subtotal decimal(14, 2),
    loaded_at datetime
)
engine = MergeTree
order by item_id;