package tenant

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tenantInterface interface {
	SetTenantID(tenantID string)
	GetTenantID() string
}

var ContextKeyTenantID = "tenant:id"
var DisableTenantScope = "tenant:disable"

func getTenantID(db *gorm.DB) (string, bool) {
	if tenantID, ok := db.Get(ContextKeyTenantID); ok {
		return fmt.Sprint(tenantID), ok
	}
	return "", false
}

func tryGetTenantModel(db *gorm.DB) (tenantInterface, bool) {
	value, enabled := db.Statement.Model.(tenantInterface)
	return value, enabled
}

func assignTenantID(db *gorm.DB) {
	if value, enabled := db.Statement.Model.(tenantInterface); enabled {
		if tenantID, ok := getTenantID(db); ok {
			value.SetTenantID(tenantID)
		}
	}
}

func addWhereCondition(db *gorm.DB) {
	if _, disabled := db.Get(DisableTenantScope); !disabled {
		if _, ok := tryGetTenantModel(db); ok {
			if tenantID, ok := getTenantID(db); ok {
				condition := db.Statement.BuildCondition("tenant_id", tenantID)
				db.Statement.AddClause(clause.Where{Exprs: condition})
			}
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

	// register update callback
	updateCallback := db.Callback().Update().Before("gorm:update")
	if err := updateCallback.Register("tenant:add_where_condition", addWhereCondition); err != nil {
		return err
	}

	// register query callback
	queryCallback := db.Callback().Query().Before("gorm:query")
	if err := queryCallback.Register("tenant:add_where_condition", addWhereCondition); err != nil {
		return err
	}

	// register delete callback
	deleteCallback := db.Callback().Delete().Before("gorm:delete")
	if err := deleteCallback.Register("tenant:add_where_condition", addWhereCondition); err != nil {
		return err
	}

	return nil
}
