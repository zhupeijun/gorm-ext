package tenant

import (
	"fmt"

	"gorm.io/gorm"
)

type tenantInterface interface {
	SetTenantID(tenantID string)
	GetTenantID() string
}

var ContextKeyTenantID = "tenant:id"

func getTenantID(db *gorm.DB) (string, bool) {
	if tenantID, ok := db.Get(ContextKeyTenantID); ok {
		return fmt.Sprint(tenantID), ok
	}
	return "", false
}

func assignTenantID(db *gorm.DB) {
	if value, enabled := db.Statement.Model.(tenantInterface); enabled {
		if tenantID, ok := getTenantID(db); ok {
			value.SetTenantID(tenantID)
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) error {
	// register create callback
	createCallback := db.Callback().Create().Before("gorm:create")
	if err := createCallback.Register("tenant:assign_tenant_id", assignTenantID); err != nil {
		return err
	}

	return nil
}
