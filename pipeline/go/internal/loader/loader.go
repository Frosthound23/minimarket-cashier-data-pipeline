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
