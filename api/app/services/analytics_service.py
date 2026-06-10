from app.databases import get_clickhouse_client


def rows_to_dicts(result):
    return [
        dict(zip(result.column_names, row))
        for row in result.result_rows
    ]


def get_revenue_by_store():
    client = get_clickhouse_client()

    query = """
        select
            fs.tenant_id,
            ds.store_name,
            ds.city,
            sum(fs.subtotal) as total_revenue,
            countDistinct(fs.transaction_key) as total_transactions
        from fact_sales fs
        left join dim_store ds
            on fs.store_key = ds.store_key
        group by
            fs.tenant_id,
            ds.store_name,
            ds.city
        order by total_revenue desc
    """

    result = client.query(query)
    return rows_to_dicts(result)


def get_promotion_effectiveness():
    client = get_clickhouse_client()

    query = """
        select
            tenant_id,
            promo_name,
            promo_type,
            countDistinct(transaction_key) as promoted_transactions,
            sum(discount_applied) as total_discount,
            sum(total_amount) as promoted_revenue
        from fact_promotion_usage
        group by
            tenant_id,
            promo_name,
            promo_type
        order by promoted_revenue desc
    """

    result = client.query(query)
    return rows_to_dicts(result)


def get_top_products_by_city():
    client = get_clickhouse_client()

    query = """
        select
            fs.tenant_id,
            ds.city,
            dp.product_name,
            dp.category,
            sum(fs.quantity) as total_quantity,
            sum(fs.subtotal) as total_revenue
        from fact_sales fs
        left join dim_store ds
            on fs.store_key = ds.store_key
        left join dim_product dp
            on fs.product_key = dp.product_key
        group by
            fs.tenant_id,
            ds.city,
            dp.product_name,
            dp.category
        order by total_revenue desc
        limit 20
    """

    result = client.query(query)
    return rows_to_dicts(result)


def get_customer_segments():
    client = get_clickhouse_client()

    query = """
        with customer_sales as (
            select
                fs.tenant_id,
                dc.customer_key,
                dc.customer_name,
                dc.city,
                countDistinct(fs.transaction_key) as total_transactions,
                sum(fs.subtotal) as total_spent
            from fact_sales fs
            left join dim_customer dc
                on fs.customer_key = dc.customer_key
            group by
                fs.tenant_id,
                dc.customer_key,
                dc.customer_name,
                dc.city
        )

        select
            tenant_id,
            case
                when total_spent >= 1000000 then 'High Value'
                when total_spent >= 500000 then 'Medium Value'
                else 'Low Value'
            end as customer_segment,
            count(*) as total_customers,
            sum(total_spent) as segment_revenue,
            avg(total_spent) as avg_customer_spent
        from customer_sales
        group by
            tenant_id,
            customer_segment
        order by
            tenant_id,
            segment_revenue desc
    """

    result = client.query(query)
    return rows_to_dicts(result)


def get_transactions_by_day():
    client = get_clickhouse_client()

    query = """
        select
            dd.day_name,
            dd.day_of_week,
            countDistinct(fs.transaction_key) as total_transactions,
            sum(fs.subtotal) as total_revenue
        from fact_sales fs
        left join dim_date dd
            on fs.transaction_date_key = dd.date_day
        group by
            dd.day_name,
            dd.day_of_week
        order by dd.day_of_week
    """

    result = client.query(query)
    return rows_to_dicts(result)


def get_summary_metrics():
    client = get_clickhouse_client()

    query = """
        select
            sum(subtotal) as total_revenue,
            countDistinct(transaction_key) as total_transactions,
            countDistinct(customer_key) as total_customers,
            sum(subtotal) / countDistinct(transaction_key) as average_basket_size
        from fact_sales
    """

    result = client.query(query)
    rows = rows_to_dicts(result)

    if not rows:
        return {
            "total_revenue": 0,
            "total_transactions": 0,
            "total_customers": 0,
            "average_basket_size": 0,
        }

    return rows[0]