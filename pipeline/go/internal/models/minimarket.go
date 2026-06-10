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
type Product struct {
	TenantID    string
	ProductID   int
	ProductName string
	Category    sql.NullString
	Brand       sql.NullString
	UnitPrice   string
	IsActive    bool
	CreatedAt   time.Time
	LoadedAt    time.Time
}

type Store struct {
	TenantID  string
	StoreID   int
	StoreName string
	City      sql.NullString
	Province  sql.NullString
	StoreType sql.NullString
	OpenedAt  sql.NullTime
	IsActive  bool
	LoadedAt  time.Time
}

type Promotion struct {
	TenantID    string
	PromoID     int
	PromoName   sql.NullString
	PromoType   sql.NullString
	DiscountPct sql.NullString
	StartDate   sql.NullTime
	EndDate     sql.NullTime
	MinPurchase string
	LoadedAt    time.Time
}

type Supplier struct {
	TenantID     string
	SupplierID   int
	SupplierName sql.NullString
	ContactName  sql.NullString
	City         sql.NullString
	Country      sql.NullString
	CreatedAt    time.Time
	LoadedAt     time.Time
}

type Transaction struct {
	TenantID        string
	TransactionID   int
	CustomerID      sql.NullInt64
	StoreID         sql.NullInt64
	TransactionDate time.Time
	TotalAmount     string
	PaymentMethod   sql.NullString
	Status          sql.NullString
	LoadedAt        time.Time
}

type TransactionItem struct {
	TenantID      string
	ItemID        int
	TransactionID sql.NullInt64
	ProductID     sql.NullInt64
	Quantity      int
	UnitPrice     string
	Discount      string
	Subtotal      string
	LoadedAt      time.Time
}

type TransactionPromotion struct {
	TenantID        string
	ID              int
	TransactionID   sql.NullInt64
	PromoID         sql.NullInt64
	DiscountApplied sql.NullString
	LoadedAt        time.Time
}
