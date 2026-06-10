package app

import (
	"context"
	"fmt"
	"log"
	"minimarket-go-pipeline/internal/extractor"
	"minimarket-go-pipeline/internal/loader"
	"minimarket-go-pipeline/internal/models"
	"sync"
	"time"
)

type Pipeline struct {
	extractor *extractor.PostgresExtractor
	loader    *loader.ClickHouseLoader
}

func NewPipeline(
	extractor *extractor.PostgresExtractor,
	loader *loader.ClickHouseLoader,
) *Pipeline {
	return &Pipeline{
		extractor: extractor,
		loader:    loader,
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

	customers, err := p.extractor.ExtractCustomers(ctx, tenant)
	if err != nil {
		return err
	}

	log.Printf("extracted customers tenant=%s rows=%d", tenant.TenantID, len(customers))

	if err := p.loader.LoadCustomers(ctx, customers); err != nil {
		return err
	}

	log.Printf(
		"loaded customers tenant=%s rows=%d duration=%s",
		tenant.TenantID,
		len(customers),
		time.Since(startedAt),
	)

	return nil
}
