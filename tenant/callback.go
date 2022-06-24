package tenant

import (
	"context"
	"reflect"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
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

// isTargetModel any type of T, *T, []T, *[]T, *[]*T, T is type of tenantInterface
func isTargetModel(db *gorm.DB) bool {
	_, ok := reflect.New(db.Statement.Schema.ModelType).Interface().(tenantInterface)
	return ok
}

// assignTenantID set tenant id to model, model need to be pointer
func assignTenantID(db *gorm.DB) {
	modelSchema, err := schema.Parse(db.Statement.Model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		_ = db.AddError(err)
		return
	}

	if tenantID, ok := getTenantID(db); ok {
		if val, enabled := db.Statement.Model.(tenantInterface); enabled {
			val.SetTenantID(tenantID)
		} else {
			// if it is pointer get the value of it
			valueOfModel := reflect.ValueOf(db.Statement.Model)
			if valueOfModel.Kind() == reflect.Pointer {
				valueOfModel = reflect.Indirect(valueOfModel)
			}

			// if it is a slice, set value for each record
			if valueOfModel.Kind() == reflect.Slice {
				for i := 0; i < valueOfModel.Len(); i++ {
					err = modelSchema.FieldsByDBName["tenant_id"].Set(context.Background(), valueOfModel.Index(i), tenantID)
					if err != nil {
						_ = db.AddError(err)
						return
					}
				}
			}
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
