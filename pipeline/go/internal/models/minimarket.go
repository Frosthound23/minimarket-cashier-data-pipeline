package models

import (
	"database/sql"
	"time"
)

type Tenant struct {
	TenantID  string `json:"tenant_id"`
	Schema    string `json:"schema"`
	StoreName string `json:"store_name"`
	City      string `json:"city"`
}

type TenantList struct {
	Tenants []Tenant `json:"tenants"`
}
type Customer struct {
	TenantID   string
	CustomerID int
	Name       string
	Phone      sql.NullString
	Email      sql.NullString
	Gender     sql.NullString
	City       sql.NullString
	CreatedAt  time.Time
	LoadedAt   time.Time
}
