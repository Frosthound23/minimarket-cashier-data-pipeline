package watermark

import (
	"context"
	"fmt"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

type Store struct {
	conn clickhouse.Conn
}

func NewStore(conn clickhouse.Conn) *Store {
	return &Store{
		conn: conn,
	}
}

func (s *Store) GetWatermark(
	ctx context.Context,
	tenantID string,
	tableName string,
) (time.Time, error) {
	query := `
		SELECT
			ifNull(max(last_watermark), toDateTime('1970-01-01 00:00:00')) AS last_watermark
		FROM elt_watermarks
		WHERE tenant_id = ?
		  AND table_name = ?
	`

	var lastWatermark time.Time

	if err := s.conn.QueryRow(
		ctx,
		query,
		tenantID,
		tableName,
	).Scan(&lastWatermark); err != nil {
		return time.Time{}, fmt.Errorf(
			"failed to get watermark tenant=%s table=%s: %w",
			tenantID,
			tableName,
			err,
		)
	}

	return lastWatermark, nil
}

func (s *Store) UpdateWatermark(
	ctx context.Context,
	tenantID string,
	tableName string,
	watermarkColumn string,
	lastWatermark time.Time,
) error {
	if lastWatermark.IsZero() {
		return nil
	}

	query := `
		INSERT INTO elt_watermarks (
			tenant_id,
			table_name,
			watermark_column,
			last_watermark,
			processed_at
		)
		VALUES (?, ?, ?, ?, ?)
	`

	if err := s.conn.Exec(
		ctx,
		query,
		tenantID,
		tableName,
		watermarkColumn,
		lastWatermark,
		time.Now().UTC(),
	); err != nil {
		return fmt.Errorf(
			"failed to update watermark tenant=%s table=%s: %w",
			tenantID,
			tableName,
			err,
		)
	}

	return nil
}
