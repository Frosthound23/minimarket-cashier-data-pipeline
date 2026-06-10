package extractor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"minimarket-go-pipeline/internal/models"
	"time"
)

type PostgresExtractor struct {
	db *sql.DB
}

func NewPostgresExtractor(db *sql.DB) *PostgresExtractor {
	return &PostgresExtractor{db: db}
}

func (e *PostgresExtractor) ExtractCustomers(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Customer, error) {
	log.Println("extracting customers for tenant", tenant.TenantID)
	query := fmt.Sprintf(`
		SELECT
			customer_id,
			name,
			phone,
			email,
			gender,
			city,
			created_at
		FROM %s.customers
		ORDER BY customer_id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract customers for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	customers := make([]models.Customer, 0)

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
			return nil, fmt.Errorf("failed to scan customer row for tenant %s: %w", tenant.TenantID, err)
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("customer row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return customers, nil
}

func (e *PostgresExtractor) ExtractProducts(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Product, error) {
	query := fmt.Sprintf(`
		SELECT
			product_id,
			product_name,
			category,
			brand,
			unit_price::text,
			is_active,
			created_at
		FROM %s.products
		ORDER BY product_id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract products for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	products := make([]models.Product, 0)

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
			return nil, fmt.Errorf("failed to scan product row for tenant %s: %w", tenant.TenantID, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("product row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return products, nil
}

func (e *PostgresExtractor) ExtractStores(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Store, error) {
	query := fmt.Sprintf(`
		SELECT
			store_id,
			store_name,
			city,
			province,
			store_type,
			opened_at,
			is_active
		FROM %s.stores
		ORDER BY store_id
	`, tenant.Schema)

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
	query := fmt.Sprintf(`
		SELECT
			promo_id,
			promo_name,
			promo_type,
			discount_pct::text,
			start_date,
			end_date,
			min_purchase::text
		FROM %s.promotions
		ORDER BY promo_id
	`, tenant.Schema)

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
) ([]models.Supplier, error) {
	query := fmt.Sprintf(`
		SELECT
			supplier_id,
			supplier_name,
			contact_name,
			city,
			country,
			created_at
		FROM %s.suppliers
		ORDER BY supplier_id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract suppliers for tenant %s: %w", tenant.TenantID, err)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	suppliers := make([]models.Supplier, 0)

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
			return nil, fmt.Errorf("failed to scan supplier row for tenant %s: %w", tenant.TenantID, err)
		}

		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("supplier row iteration failed for tenant %s: %w", tenant.TenantID, err)
	}

	return suppliers, nil
}

func (e *PostgresExtractor) ExtractTransactions(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Transaction, error) {
	query := fmt.Sprintf(`
		SELECT
			transaction_id,
			customer_id,
			store_id,
			transaction_date,
			total_amount::text,
			payment_method,
			status
		FROM %s.transactions
		ORDER BY transaction_id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to extract transactions for tenant %s: %w",
			tenant.TenantID,
			err,
		)
	}
	defer rows.Close()

	loadedAt := time.Now().UTC()
	transactions := make([]models.Transaction, 0)

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
			return nil, fmt.Errorf(
				"failed to scan transaction row for tenant %s: %w",
				tenant.TenantID,
				err,
			)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"transaction row iteration failed for tenant %s: %w",
			tenant.TenantID,
			err,
		)
	}

	return transactions, nil
}

func (e *PostgresExtractor) ExtractTransactionItems(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.TransactionItem, error) {
	query := fmt.Sprintf(`
		SELECT
			item_id,
			transaction_id,
			product_id,
			quantity,
			unit_price::text,
			discount::text,
			subtotal::text
		FROM %s.transaction_items
		ORDER BY item_id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to extract transaction_items for tenant %s: %w",
			tenant.TenantID,
			err,
		)
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
			return nil, fmt.Errorf(
				"failed to scan transaction_item row for tenant %s: %w",
				tenant.TenantID,
				err,
			)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"transaction_item row iteration failed for tenant %s: %w",
			tenant.TenantID,
			err,
		)
	}

	return items, nil
}

func (e *PostgresExtractor) ExtractTransactionPromotions(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.TransactionPromotion, error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			transaction_id,
			promo_id,
			discount_applied::text
		FROM %s.transaction_promotions
		ORDER BY id
	`, tenant.Schema)

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to extract transaction_promotions for tenant %s: %w",
			tenant.TenantID,
			err,
		)
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
			return nil, fmt.Errorf(
				"failed to scan transaction_promotion row for tenant %s: %w",
				tenant.TenantID,
				err,
			)
		}

		promotions = append(promotions, promotion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"transaction_promotion row iteration failed for tenant %s: %w",
			tenant.TenantID,
			err,
		)
	}

	return promotions, nil
}
