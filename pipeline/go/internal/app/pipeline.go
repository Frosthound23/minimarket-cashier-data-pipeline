package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"minimarket-go-pipeline/internal/extractor"
	"minimarket-go-pipeline/internal/loader"
	"minimarket-go-pipeline/internal/models"
	"minimarket-go-pipeline/internal/watermark"
)

type Pipeline struct {
	extractor      *extractor.PostgresExtractor
	loader         *loader.ClickHouseLoader
	watermarkStore *watermark.Store
}

func NewPipeline(
	extractor *extractor.PostgresExtractor,
	loader *loader.ClickHouseLoader,
	watermarkStore *watermark.Store,
) *Pipeline {
	return &Pipeline{
		extractor:      extractor,
		loader:         loader,
		watermarkStore: watermarkStore,
	}
}

func (p *Pipeline) Run(ctx context.Context, tenants []models.Tenant) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(tenants))

	for _, tenant := range tenants {
		wg.Add(1)

		go func(t models.Tenant) {
			defer wg.Done()

			if err := p.runTenant(ctx, t); err != nil {
				errCh <- fmt.Errorf("tenant %s failed: %w", t.TenantID, err)
			}
		}(tenant)
	}

	wg.Wait()
	close(errCh)

	hasError := false

	for err := range errCh {
		hasError = true
		log.Println(err)
	}

	if hasError {
		return fmt.Errorf("one or more tenant pipelines failed")
	}

	return nil
}

func (p *Pipeline) runTenant(ctx context.Context, tenant models.Tenant) error {
	startedAt := time.Now()

	log.Printf("starting tenant=%s schema=%s", tenant.TenantID, tenant.Schema)

	if err := p.loadIncrementalCustomers(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadIncrementalProducts(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadIncrementalSuppliers(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadIncrementalTransactions(ctx, tenant); err != nil {
		return err
	}

	if err := p.clearFullRefreshTables(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadFullRefreshStores(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadFullRefreshPromotions(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadFullRefreshTransactionItems(ctx, tenant); err != nil {
		return err
	}

	if err := p.loadFullRefreshTransactionPromotions(ctx, tenant); err != nil {
		return err
	}

	log.Printf(
		"finished tenant=%s duration=%s",
		tenant.TenantID,
		time.Since(startedAt),
	)

	return nil
}

func (p *Pipeline) loadIncrementalCustomers(
	ctx context.Context,
	tenant models.Tenant,
) error {
	lastWatermark, err := p.watermarkStore.GetWatermark(
		ctx,
		tenant.TenantID,
		"customers",
	)
	if err != nil {
		return err
	}

	customers, newWatermark, err := p.extractor.ExtractCustomers(
		ctx,
		tenant,
		lastWatermark,
	)
	if err != nil {
		return err
	}

	if err := p.loader.LoadCustomers(ctx, customers); err != nil {
		return err
	}

	if len(customers) > 0 {
		if err := p.watermarkStore.UpdateWatermark(
			ctx,
			tenant.TenantID,
			"customers",
			"created_at",
			newWatermark,
		); err != nil {
			return err
		}
	}

	log.Printf(
		"loaded customers tenant=%s rows=%d previous_watermark=%s new_watermark=%s",
		tenant.TenantID,
		len(customers),
		lastWatermark.Format(time.RFC3339),
		newWatermark.Format(time.RFC3339),
	)

	return nil
}

func (p *Pipeline) loadIncrementalProducts(
	ctx context.Context,
	tenant models.Tenant,
) error {
	lastWatermark, err := p.watermarkStore.GetWatermark(
		ctx,
		tenant.TenantID,
		"products",
	)
	if err != nil {
		return err
	}

	products, newWatermark, err := p.extractor.ExtractProducts(
		ctx,
		tenant,
		lastWatermark,
	)
	if err != nil {
		return err
	}

	if err := p.loader.LoadProducts(ctx, products); err != nil {
		return err
	}

	if len(products) > 0 {
		if err := p.watermarkStore.UpdateWatermark(
			ctx,
			tenant.TenantID,
			"products",
			"created_at",
			newWatermark,
		); err != nil {
			return err
		}
	}

	log.Printf(
		"loaded products tenant=%s rows=%d previous_watermark=%s new_watermark=%s",
		tenant.TenantID,
		len(products),
		lastWatermark.Format(time.RFC3339),
		newWatermark.Format(time.RFC3339),
	)

	return nil
}

func (p *Pipeline) loadIncrementalSuppliers(
	ctx context.Context,
	tenant models.Tenant,
) error {
	lastWatermark, err := p.watermarkStore.GetWatermark(
		ctx,
		tenant.TenantID,
		"suppliers",
	)
	if err != nil {
		return err
	}

	suppliers, newWatermark, err := p.extractor.ExtractSuppliers(
		ctx,
		tenant,
		lastWatermark,
	)
	if err != nil {
		return err
	}

	if err := p.loader.LoadSuppliers(ctx, suppliers); err != nil {
		return err
	}

	if len(suppliers) > 0 {
		if err := p.watermarkStore.UpdateWatermark(
			ctx,
			tenant.TenantID,
			"suppliers",
			"created_at",
			newWatermark,
		); err != nil {
			return err
		}
	}

	log.Printf(
		"loaded suppliers tenant=%s rows=%d previous_watermark=%s new_watermark=%s",
		tenant.TenantID,
		len(suppliers),
		lastWatermark.Format(time.RFC3339),
		newWatermark.Format(time.RFC3339),
	)

	return nil
}

func (p *Pipeline) loadIncrementalTransactions(
	ctx context.Context,
	tenant models.Tenant,
) error {
	lastWatermark, err := p.watermarkStore.GetWatermark(
		ctx,
		tenant.TenantID,
		"transactions",
	)
	if err != nil {
		return err
	}

	transactions, newWatermark, err := p.extractor.ExtractTransactions(
		ctx,
		tenant,
		lastWatermark,
	)
	if err != nil {
		return err
	}

	if err := p.loader.LoadTransactions(ctx, transactions); err != nil {
		return err
	}

	if len(transactions) > 0 {
		if err := p.watermarkStore.UpdateWatermark(
			ctx,
			tenant.TenantID,
			"transactions",
			"transaction_date",
			newWatermark,
		); err != nil {
			return err
		}
	}

	log.Printf(
		"loaded transactions tenant=%s rows=%d previous_watermark=%s new_watermark=%s",
		tenant.TenantID,
		len(transactions),
		lastWatermark.Format(time.RFC3339),
		newWatermark.Format(time.RFC3339),
	)

	return nil
}

func (p *Pipeline) clearFullRefreshTables(
	ctx context.Context,
	tenant models.Tenant,
) error {
	fullRefreshTables := []string{
		"raw_stores",
		"raw_promotions",
		"raw_transaction_items",
		"raw_transaction_promotions",
	}

	if err := p.loader.ClearTenantTables(
		ctx,
		tenant.TenantID,
		fullRefreshTables,
	); err != nil {
		return err
	}

	log.Printf("cleared full-refresh tables tenant=%s", tenant.TenantID)

	return nil
}

func (p *Pipeline) loadFullRefreshStores(
	ctx context.Context,
	tenant models.Tenant,
) error {
	stores, err := p.extractor.ExtractStores(ctx, tenant)
	if err != nil {
		return err
	}

	if err := p.loader.LoadStores(ctx, stores); err != nil {
		return err
	}

	log.Printf("loaded stores tenant=%s rows=%d", tenant.TenantID, len(stores))

	return nil
}

func (p *Pipeline) loadFullRefreshPromotions(
	ctx context.Context,
	tenant models.Tenant,
) error {
	promotions, err := p.extractor.ExtractPromotions(ctx, tenant)
	if err != nil {
		return err
	}

	if err := p.loader.LoadPromotions(ctx, promotions); err != nil {
		return err
	}

	log.Printf("loaded promotions tenant=%s rows=%d", tenant.TenantID, len(promotions))

	return nil
}

func (p *Pipeline) loadFullRefreshTransactionItems(
	ctx context.Context,
	tenant models.Tenant,
) error {
	transactionItems, err := p.extractor.ExtractTransactionItems(ctx, tenant)
	if err != nil {
		return err
	}

	if err := p.loader.LoadTransactionItems(ctx, transactionItems); err != nil {
		return err
	}

	log.Printf(
		"loaded transaction_items tenant=%s rows=%d",
		tenant.TenantID,
		len(transactionItems),
	)

	return nil
}

func (p *Pipeline) loadFullRefreshTransactionPromotions(
	ctx context.Context,
	tenant models.Tenant,
) error {
	transactionPromotions, err := p.extractor.ExtractTransactionPromotions(ctx, tenant)
	if err != nil {
		return err
	}

	if err := p.loader.LoadTransactionPromotions(ctx, transactionPromotions); err != nil {
		return err
	}

	log.Printf(
		"loaded transaction_promotions tenant=%s rows=%d",
		tenant.TenantID,
		len(transactionPromotions),
	)

	return nil
}
