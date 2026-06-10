from fastapi import APIRouter

from app.services.analytics_service import (
    get_customer_segments,
    get_promotion_effectiveness,
    get_revenue_by_store,
    get_summary_metrics,
    get_top_products_by_city,
    get_transactions_by_day,
)

router = APIRouter(prefix="/api", tags=["analytics"])


@router.get("/summary")
def summary():
    return get_summary_metrics()


@router.get("/revenue-by-store")
def revenue_by_store():
    return get_revenue_by_store()


@router.get("/promotion-effectiveness")
def promotion_effectiveness():
    return get_promotion_effectiveness()


@router.get("/top-products-by-city")
def top_products_by_city():
    return get_top_products_by_city()


@router.get("/customer-segments")
def customer_segments():
    return get_customer_segments()


@router.get("/transactions-by-day")
def transactions_by_day():
    return get_transactions_by_day()