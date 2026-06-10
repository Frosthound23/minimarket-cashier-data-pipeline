package loader

import (
	"context"
	"fmt"
	"minimarket-go-pipeline/internal/helper"
	"minimarket-go-pipeline/internal/models"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseLoader struct {
	conn clickhouse.Conn
}

func NewClickHouseLoader(conn clickhouse.Conn) *ClickHouseLoader {
	return &ClickHouseLoader{
		conn: conn,
	}
}

func (l *ClickHouseLoader) LoadCustomers(
	ctx context.Context,
	customers []models.Customer,
) error {
	if len(customers) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(
		ctx,
		`
		INSERT INTO raw_customers (
			tenant_id,
			customer_id,
			name,
			phone,
			email,
			gender,
			city,
			created_at,
			loaded_at
		)
		`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_customers batch: %w", err)
	}

	for _, customer := range customers {
		if err := batch.Append(
			customer.TenantID,
			int32(customer.CustomerID),
			customer.Name,
			helper.NullableStringSQL(customer.Phone),
			helper.NullableStringSQL(customer.Email),
			helper.NullableStringSQL(customer.Gender),
			helper.NullableStringSQL(customer.City),
			customer.CreatedAt,
			customer.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append customer row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_customers batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadProducts(
	ctx context.Context,
	products []models.Product,
) error {
	if len(products) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_products (
			tenant_id,
			product_id,
			product_name,
			category,
			brand,
			unit_price,
			is_active,
			created_at,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_products batch: %w", err)
	}

	for _, product := range products {
		if err := batch.Append(
			product.TenantID,
			int32(product.ProductID),
			product.ProductName,
			helper.NullableStringSQL(product.Category),
			helper.NullableStringSQL(product.Brand),
			product.UnitPrice,
			product.IsActive,
			product.CreatedAt,
			product.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append product row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_products batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadStores(
	ctx context.Context,
	stores []models.Store,
) error {
	if len(stores) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_stores (
			tenant_id,
			store_id,
			store_name,
			city,
			province,
			store_type,
			opened_at,
			is_active,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_stores batch: %w", err)
	}

	for _, store := range stores {
		if err := batch.Append(
			store.TenantID,
			int32(store.StoreID),
			store.StoreName,
			helper.NullableStringSQL(store.City),
			helper.NullableStringSQL(store.Province),
			helper.NullableStringSQL(store.StoreType),
			helper.NullableTimeSQL(store.OpenedAt),
			store.IsActive,
			store.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append store row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_stores batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadPromotions(
	ctx context.Context,
	promotions []models.Promotion,
) error {
	if len(promotions) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_promotions (
			tenant_id,
			promo_id,
			promo_name,
			promo_type,
			discount_pct,
			start_date,
			end_date,
			min_purchase,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_promotions batch: %w", err)
	}

	for _, promotion := range promotions {
		if err := batch.Append(
			promotion.TenantID,
			int32(promotion.PromoID),
			helper.NullableStringSQL(promotion.PromoName),
			helper.NullableStringSQL(promotion.PromoType),
			helper.NullableStringSQL(promotion.DiscountPct),
			helper.NullableTimeSQL(promotion.StartDate),
			helper.NullableTimeSQL(promotion.EndDate),
			promotion.MinPurchase,
			promotion.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append promotion row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_promotions batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadSuppliers(
	ctx context.Context,
	suppliers []models.Supplier,
) error {
	if len(suppliers) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_suppliers (
			tenant_id,
			supplier_id,
			supplier_name,
			contact_name,
			city,
			country,
			created_at,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_suppliers batch: %w", err)
	}

	for _, supplier := range suppliers {
		if err := batch.Append(
			supplier.TenantID,
			int32(supplier.SupplierID),
			helper.NullableStringSQL(supplier.SupplierName),
			helper.NullableStringSQL(supplier.ContactName),
			helper.NullableStringSQL(supplier.City),
			helper.NullableStringSQL(supplier.Country),
			supplier.CreatedAt,
			supplier.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append supplier row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_suppliers batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadTransactions(
	ctx context.Context,
	transactions []models.Transaction,
) error {
	if len(transactions) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_transactions (
			tenant_id,
			transaction_id,
			customer_id,
			store_id,
			transaction_date,
			total_amount,
			payment_method,
			status,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_transactions batch: %w", err)
	}

	for _, transaction := range transactions {
		if err := batch.Append(
			transaction.TenantID,
			int32(transaction.TransactionID),
			helper.NullableInt32SQL(transaction.CustomerID),
			helper.NullableInt32SQL(transaction.StoreID),
			transaction.TransactionDate,
			transaction.TotalAmount,
			helper.NullableStringSQL(transaction.PaymentMethod),
			helper.NullableStringSQL(transaction.Status),
			transaction.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_transactions batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadTransactionItems(
	ctx context.Context,
	items []models.TransactionItem,
) error {
	if len(items) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_transaction_items (
			tenant_id,
			item_id,
			transaction_id,
			product_id,
			quantity,
			unit_price,
			discount,
			subtotal,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_transaction_items batch: %w", err)
	}

	for _, item := range items {
		if err := batch.Append(
			item.TenantID,
			int32(item.ItemID),
			helper.NullableInt32SQL(item.TransactionID),
			helper.NullableInt32SQL(item.ProductID),
			int32(item.Quantity),
			item.UnitPrice,
			item.Discount,
			item.Subtotal,
			item.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction_item row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_transaction_items batch: %w", err)
	}

	return nil
}

func (l *ClickHouseLoader) LoadTransactionPromotions(
	ctx context.Context,
	promotions []models.TransactionPromotion,
) error {
	if len(promotions) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
		INSERT INTO raw_transaction_promotions (
			tenant_id,
			id,
			transaction_id,
			promo_id,
			discount_applied,
			loaded_at
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare raw_transaction_promotions batch: %w", err)
	}

	for _, promotion := range promotions {
		if err := batch.Append(
			promotion.TenantID,
			int32(promotion.ID),
			helper.NullableInt32SQL(promotion.TransactionID),
			helper.NullableInt32SQL(promotion.PromoID),
			helper.NullableStringSQL(promotion.DiscountApplied),
			promotion.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction_promotion row: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_transaction_promotions batch: %w", err)
	}

	return nil
}
