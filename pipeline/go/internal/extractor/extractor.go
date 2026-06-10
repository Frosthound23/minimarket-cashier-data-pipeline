package extractor

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/lib/pq"

	"minimarket-go-pipeline/internal/models"
)

var validSchemaName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type PostgresExtractor struct {
	db *sql.DB
}

func NewPostgresExtractor(db *sql.DB) *PostgresExtractor {
	return &PostgresExtractor{db: db}
}

func qualifiedTenantTable(tenant models.Tenant, tableName string) (string, error) {
	if !validSchemaName.MatchString(tenant.Schema) {
		return "", fmt.Errorf("invalid schema name for tenant %s: %s", tenant.TenantID, tenant.Schema)
	}

	return fmt.Sprintf("%s.%s", pq.QuoteIdentifier(tenant.Schema), pq.QuoteIdentifier(tableName)), nil
}

func (e *PostgresExtractor) ExtractCustomers(
	ctx context.Context,
	tenant models.Tenant,
	lastWatermark time.Time,
) ([]models.Customer, time.Time, error) {
	tableName, err := qualifiedTenantTable(tenant, "customers")
	if err != nil {
		return nil, time.Time{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			customer_id,
			name,
			phone,
			email,
			gender,
			city,
			created_at
		FROM %s
		WHERE created_at > $1
		ORDER BY created_at, customer_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query, lastWatermark)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to extract customers for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	customers := make([]models.Customer, 0)
	maxWatermark := lastWatermark

	for rows.Next() {
		var customer models.Customer

		customer.TenantID = tenant.TenantID
		customer.LoadedAt = loadedAt

		if err := rows.Scan(
			&customer.CustomerID,
			&customer.Name,
			&customer.Phone,
			&customer.Email,
			&customer.Gender,
			&customer.City,
			&customer.CreatedAt,
		); err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to scan customer row for tenant %s: %w", tenant.TenantID, err)
		}

		if customer.CreatedAt.After(maxWatermark) {
			maxWatermark = customer.CreatedAt
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, fmt.Errorf("customer row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return customers, maxWatermark, nil
}

func (e *PostgresExtractor) ExtractProducts(
	ctx context.Context,
	tenant models.Tenant,
	lastWatermark time.Time,
) ([]models.Product, time.Time, error) {
	tableName, err := qualifiedTenantTable(tenant, "products")
	if err != nil {
		return nil, time.Time{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			product_id,
			product_name,
			category,
			brand,
			unit_price::text,
			is_active,
			created_at
		FROM %s
		WHERE created_at > $1
		ORDER BY created_at, product_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query, lastWatermark)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to extract products for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	products := make([]models.Product, 0)
	maxWatermark := lastWatermark

	for rows.Next() {
		var product models.Product

		product.TenantID = tenant.TenantID
		product.LoadedAt = loadedAt

		if err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Category,
			&product.Brand,
			&product.UnitPrice,
			&product.IsActive,
			&product.CreatedAt,
		); err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to scan product row for tenant %s: %w", tenant.TenantID, err)
		}

		if product.CreatedAt.After(maxWatermark) {
			maxWatermark = product.CreatedAt
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, fmt.Errorf("product row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return products, maxWatermark, nil
}

func (e *PostgresExtractor) ExtractStores(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Store, error) {
	tableName, err := qualifiedTenantTable(tenant, "stores")
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT
			store_id,
			store_name,
			city,
			province,
			store_type,
			opened_at,
			is_active
		FROM %s
		ORDER BY store_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract stores for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	stores := make([]models.Store, 0)

	for rows.Next() {
		var store models.Store

		store.TenantID = tenant.TenantID
		store.LoadedAt = loadedAt

		if err := rows.Scan(
			&store.StoreID,
			&store.StoreName,
			&store.City,
			&store.Province,
			&store.StoreType,
			&store.OpenedAt,
			&store.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to scan store row for tenant %s: %w", tenant.TenantID, err)
		}

		stores = append(stores, store)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return stores, nil
}

func (e *PostgresExtractor) ExtractPromotions(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Promotion, error) {
	tableName, err := qualifiedTenantTable(tenant, "promotions")
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT
			promo_id,
			promo_name,
			promo_type,
			discount_pct::text,
			start_date,
			end_date,
			min_purchase::text
		FROM %s
		ORDER BY promo_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract promotions for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	promotions := make([]models.Promotion, 0)

	for rows.Next() {
		var promotion models.Promotion

		promotion.TenantID = tenant.TenantID
		promotion.LoadedAt = loadedAt

		if err := rows.Scan(
			&promotion.PromoID,
			&promotion.PromoName,
			&promotion.PromoType,
			&promotion.DiscountPct,
			&promotion.StartDate,
			&promotion.EndDate,
			&promotion.MinPurchase,
		); err != nil {
			return nil, fmt.Errorf("failed to scan promotion row for tenant %s: %w", tenant.TenantID, err)
		}

		promotions = append(promotions, promotion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("promotion row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return promotions, nil
}

func (e *PostgresExtractor) ExtractSuppliers(
	ctx context.Context,
	tenant models.Tenant,
	lastWatermark time.Time,
) ([]models.Supplier, time.Time, error) {
	tableName, err := qualifiedTenantTable(tenant, "suppliers")
	if err != nil {
		return nil, time.Time{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			supplier_id,
			supplier_name,
			contact_name,
			city,
			country,
			created_at
		FROM %s
		WHERE created_at > $1
		ORDER BY created_at, supplier_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query, lastWatermark)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to extract suppliers for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	suppliers := make([]models.Supplier, 0)
	maxWatermark := lastWatermark

	for rows.Next() {
		var supplier models.Supplier

		supplier.TenantID = tenant.TenantID
		supplier.LoadedAt = loadedAt

		if err := rows.Scan(
			&supplier.SupplierID,
			&supplier.SupplierName,
			&supplier.ContactName,
			&supplier.City,
			&supplier.Country,
			&supplier.CreatedAt,
		); err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to scan supplier row for tenant %s: %w", tenant.TenantID, err)
		}

		if supplier.CreatedAt.After(maxWatermark) {
			maxWatermark = supplier.CreatedAt
		}

		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, fmt.Errorf("supplier row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return suppliers, maxWatermark, nil
}

func (e *PostgresExtractor) ExtractTransactions(
	ctx context.Context,
	tenant models.Tenant,
	lastWatermark time.Time,
) ([]models.Transaction, time.Time, error) {
	tableName, err := qualifiedTenantTable(tenant, "transactions")
	if err != nil {
		return nil, time.Time{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			transaction_id,
			customer_id,
			store_id,
			transaction_date,
			total_amount::text,
			payment_method,
			status
		FROM %s
		WHERE transaction_date > $1
		ORDER BY transaction_date, transaction_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query, lastWatermark)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to extract transactions for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	transactions := make([]models.Transaction, 0)
	maxWatermark := lastWatermark

	for rows.Next() {
		var transaction models.Transaction

		transaction.TenantID = tenant.TenantID
		transaction.LoadedAt = loadedAt

		if err := rows.Scan(
			&transaction.TransactionID,
			&transaction.CustomerID,
			&transaction.StoreID,
			&transaction.TransactionDate,
			&transaction.TotalAmount,
			&transaction.PaymentMethod,
			&transaction.Status,
		); err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to scan transaction row for tenant %s: %w", tenant.TenantID, err)
		}

		if transaction.TransactionDate.After(maxWatermark) {
			maxWatermark = transaction.TransactionDate
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, time.Time{}, fmt.Errorf("transaction row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return transactions, maxWatermark, nil
}

func (e *PostgresExtractor) ExtractTransactionItems(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.TransactionItem, error) {
	tableName, err := qualifiedTenantTable(tenant, "transaction_items")
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT
			item_id,
			transaction_id,
			product_id,
			quantity,
			unit_price::text,
			discount::text,
			subtotal::text
		FROM %s
		ORDER BY item_id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transaction_items for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	items := make([]models.TransactionItem, 0)

	for rows.Next() {
		var item models.TransactionItem

		item.TenantID = tenant.TenantID
		item.LoadedAt = loadedAt

		if err := rows.Scan(
			&item.ItemID,
			&item.TransactionID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
			&item.Subtotal,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction_item row for tenant %s: %w", tenant.TenantID, err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transaction_item row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return items, nil
}

func (e *PostgresExtractor) ExtractTransactionPromotions(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.TransactionPromotion, error) {
	tableName, err := qualifiedTenantTable(tenant, "transaction_promotions")
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			transaction_id,
			promo_id,
			discount_applied::text
		FROM %s
		ORDER BY id
	`, tableName)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transaction_promotions for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	promotions := make([]models.TransactionPromotion, 0)

	for rows.Next() {
		var promotion models.TransactionPromotion

		promotion.TenantID = tenant.TenantID
		promotion.LoadedAt = loadedAt

		if err := rows.Scan(
			&promotion.ID,
			&promotion.TransactionID,
			&promotion.PromoID,
			&promotion.DiscountApplied,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction_promotion row for tenant %s: %w", tenant.TenantID, err)
		}

		promotions = append(promotions, promotion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transaction_promotion row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return promotions, nil
}
