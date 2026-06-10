package config

import (
	"encoding/json"
	"fmt"
	"minimarket-pipeline/internal/models"
	"os"
)

func LoadTenants(path string) ([]models.Tenant, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant config: %w", err)
	}
	defer file.Close()

	var tenants []models.Tenant

	if err := json.NewDecoder(file).Decode(&tenants); err != nil {
		return nil, fmt.Errorf("failed to decode tenant config: %w", err)
	}

	if len(tenants) == 0 {
		return nil, fmt.Errorf("tenant config is empty")
	}

	for _, tenant := range tenants {
		if tenant.TenantID == "" {
			return nil, fmt.Errorf("tenant_id is required")
		}

		if tenant.Schema == "" {
			return nil, fmt.Errorf("schema is required for tenant %s", tenant.TenantID)
		}

		if tenant.StoreName == "" {
			return nil, fmt.Errorf("store_name is required for tenant %s", tenant.TenantID)
		}
	}

	return tenants, nil
}
