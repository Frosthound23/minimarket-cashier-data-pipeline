package extract

import (
	"context"
	"database/sql"
	"fmt"
	"minimarket-go-pipeline/internal/models"
	"time"
)

type PostrgresExtract struct {
	db *sql.DB
}

func NewPostgresExtractor(db *sql.DB) *PostrgresExtract {
	return &PostrgresExtract{
		db: db,
	}
}

func (e *PostrgresExtract) ExtractCustomers(
	ctx context.Context,
	tenant models.Tenant,
) ([]models.Customer, error) {
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
		return nil, fmt.Errorf(
			"failed to extract customers for tenant %s: %w",
			tenant.TenantID,
			err,
		)
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
			return nil, fmt.Errorf(
				"failed to scan customer row for tenant %s: %w",
				tenant.TenantID,
				err,
			)
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"customer rows iteration failed for tenant %s: %w",
			tenant.TenantID,
			err,
		)
	}

	return customers, nil
}
