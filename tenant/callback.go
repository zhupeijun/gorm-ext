package tenant

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tenantInterface interface {
	SetTenantID(tenantID uint)
}

var ContextKeyTenantID = "tenant:id"
var DisableTenantScope = "tenant:disable"

// getTenantID retrieve tenant id from context
func getTenantID(db *gorm.DB) (uint, bool) {
	if tenantID, ok := db.Get(ContextKeyTenantID); ok {
		return tenantID.(uint), ok
	}
	return 0, false
}

// toPointer make a struct to pointer type
func toPointer(v interface{}) interface{} {
	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		p := reflect.New(reflect.TypeOf(v))
		v = p.Interface()
	}
	return v
}

// isTargetModel any type of T, *T, []T, *[]T, *[]*T, T is type of tenantInterface
func isTargetModel(db *gorm.DB) bool {
	_, ok := reflect.New(db.Statement.Schema.ModelType).Interface().(tenantInterface)
	return ok
}

// assignTenantID set tenant id to model, model need to be pointer
func assignTenantID(db *gorm.DB) {
	if val, enabled := db.Statement.Model.(tenantInterface); enabled {
		if tenantID, ok := getTenantID(db); ok {
			val.SetTenantID(tenantID)
		}
	}
}

// addWhereCondition append where condition
func addWhereCondition(db *gorm.DB) {
	if _, disabled := db.Get(DisableTenantScope); !disabled {
		if ok := isTargetModel(db); ok {
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
