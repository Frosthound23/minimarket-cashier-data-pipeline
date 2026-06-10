package loader

import (
	"context"
	"database/sql"
	"fmt"

	"minimarket-go-pipeline/internal/helper"
	"minimarket-go-pipeline/internal/models"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/shopspring/decimal"
)

var allowedRawTables = map[string]struct{}{
	"raw_customers":              {},
	"raw_products":               {},
	"raw_stores":                 {},
	"raw_promotions":             {},
	"raw_suppliers":              {},
	"raw_transactions":           {},
	"raw_transaction_items":      {},
	"raw_transaction_promotions": {},
}

type ClickHouseLoader struct {
	conn clickhouse.Conn
}

func NewClickHouseLoader(conn clickhouse.Conn) *ClickHouseLoader {
	return &ClickHouseLoader{
		conn: conn,
	}
}

func (l *ClickHouseLoader) ClearTenantTables(
	ctx context.Context,
	tenantID string,
	tables []string,
) error {
	for _, table := range tables {
		if _, ok := allowedRawTables[table]; !ok {
			return fmt.Errorf("table %s is not allowed to be cleared", table)
		}

		query := fmt.Sprintf(
			"ALTER TABLE %s DELETE WHERE tenant_id = ?",
			table,
		)

		if err := l.conn.Exec(ctx, query, tenantID); err != nil {
			return fmt.Errorf(
				"failed to clear table %s for tenant %s: %w",
				table,
				tenantID,
				err,
			)
		}
	}

	return nil
}

func (l *ClickHouseLoader) ClearFullRefreshTables(
	ctx context.Context,
	tenantID string,
) error {
	fullRefreshTables := []string{
		"raw_stores",
		"raw_promotions",
		"raw_transaction_items",
		"raw_transaction_promotions",
	}

	return l.ClearTenantTables(ctx, tenantID, fullRefreshTables)
}

func (l *ClickHouseLoader) LoadCustomers(
	ctx context.Context,
	customers []models.Customer,
) error {
	if len(customers) == 0 {
		return nil
	}

	batch, err := l.conn.PrepareBatch(ctx, `
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
	`)
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
			return fmt.Errorf("failed to append customer row tenant=%s customer_id=%d: %w", customer.TenantID, customer.CustomerID, err)
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
		unitPrice, err := decimalFromString(product.UnitPrice)
		if err != nil {
			return fmt.Errorf("invalid unit_price tenant=%s product_id=%d: %w", product.TenantID, product.ProductID, err)
		}

		if err := batch.Append(
			product.TenantID,
			int32(product.ProductID),
			product.ProductName,
			helper.NullableStringSQL(product.Category),
			helper.NullableStringSQL(product.Brand),
			unitPrice,
			product.IsActive,
			product.CreatedAt,
			product.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append product row tenant=%s product_id=%d: %w", product.TenantID, product.ProductID, err)
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
			return fmt.Errorf("failed to append store row tenant=%s store_id=%d: %w", store.TenantID, store.StoreID, err)
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
		discountPct, err := nullableDecimalFromSQLString(promotion.DiscountPct)
		if err != nil {
			return fmt.Errorf("invalid discount_pct tenant=%s promo_id=%d: %w", promotion.TenantID, promotion.PromoID, err)
		}

		minPurchase, err := decimalFromString(promotion.MinPurchase)
		if err != nil {
			return fmt.Errorf("invalid min_purchase tenant=%s promo_id=%d: %w", promotion.TenantID, promotion.PromoID, err)
		}

		if err := batch.Append(
			promotion.TenantID,
			int32(promotion.PromoID),
			helper.NullableStringSQL(promotion.PromoName),
			helper.NullableStringSQL(promotion.PromoType),
			discountPct,
			helper.NullableTimeSQL(promotion.StartDate),
			helper.NullableTimeSQL(promotion.EndDate),
			minPurchase,
			promotion.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append promotion row tenant=%s promo_id=%d: %w", promotion.TenantID, promotion.PromoID, err)
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
			return fmt.Errorf("failed to append supplier row tenant=%s supplier_id=%d: %w", supplier.TenantID, supplier.SupplierID, err)
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
		totalAmount, err := decimalFromString(transaction.TotalAmount)
		if err != nil {
			return fmt.Errorf("invalid total_amount tenant=%s transaction_id=%d: %w", transaction.TenantID, transaction.TransactionID, err)
		}

		if err := batch.Append(
			transaction.TenantID,
			int32(transaction.TransactionID),
			helper.NullableInt32SQL(transaction.CustomerID),
			helper.NullableInt32SQL(transaction.StoreID),
			transaction.TransactionDate,
			totalAmount,
			helper.NullableStringSQL(transaction.PaymentMethod),
			helper.NullableStringSQL(transaction.Status),
			transaction.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction row tenant=%s transaction_id=%d: %w", transaction.TenantID, transaction.TransactionID, err)
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
		unitPrice, err := decimalFromString(item.UnitPrice)
		if err != nil {
			return fmt.Errorf("invalid unit_price tenant=%s item_id=%d: %w", item.TenantID, item.ItemID, err)
		}

		discount, err := decimalFromString(item.Discount)
		if err != nil {
			return fmt.Errorf("invalid discount tenant=%s item_id=%d: %w", item.TenantID, item.ItemID, err)
		}

		subtotal, err := decimalFromString(item.Subtotal)
		if err != nil {
			return fmt.Errorf("invalid subtotal tenant=%s item_id=%d: %w", item.TenantID, item.ItemID, err)
		}

		if err := batch.Append(
			item.TenantID,
			int32(item.ItemID),
			helper.NullableInt32SQL(item.TransactionID),
			helper.NullableInt32SQL(item.ProductID),
			int32(item.Quantity),
			unitPrice,
			discount,
			subtotal,
			item.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction_item row tenant=%s item_id=%d: %w", item.TenantID, item.ItemID, err)
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
		discountApplied, err := nullableDecimalFromSQLString(promotion.DiscountApplied)
		if err != nil {
			return fmt.Errorf("invalid discount_applied tenant=%s id=%d: %w", promotion.TenantID, promotion.ID, err)
		}

		if err := batch.Append(
			promotion.TenantID,
			int32(promotion.ID),
			helper.NullableInt32SQL(promotion.TransactionID),
			helper.NullableInt32SQL(promotion.PromoID),
			discountApplied,
			promotion.LoadedAt,
		); err != nil {
			return fmt.Errorf("failed to append transaction_promotion row tenant=%s id=%d: %w", promotion.TenantID, promotion.ID, err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send raw_transaction_promotions batch: %w", err)
	}

	return nil
}

func decimalFromString(value string) (decimal.Decimal, error) {
	if value == "" {
		return decimal.Zero, nil
	}

	parsedValue, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return parsedValue, nil
}

func nullableDecimalFromSQLString(value sql.NullString) (*decimal.Decimal, error) {
	if !value.Valid {
		return nil, nil
	}

	parsedValue, err := decimal.NewFromString(value.String)
	if err != nil {
		return nil, err
	}

	return &parsedValue, nil
}
