package audit

import (
	"fmt"

	"gorm.io/gorm"
)

type auditInterface interface {
	SetCreatedBy(createdBy string)
	GetCreatedBy() string
	SetUpdatedBy(updatedBy string)
	GetUpdatedBy() string
}

var ContextKeyCurrentUser = "audit:current_user"

func getCurrentUser(db *gorm.DB) (string, bool) {
	if user, ok := db.Get(ContextKeyCurrentUser); ok {
		return fmt.Sprint(user), ok
	}
	return "", false
}

func assignCreatedBy(db *gorm.DB) {
	if value, enabled := db.Statement.Model.(auditInterface); enabled {
		if user, ok := getCurrentUser(db); ok {
			value.SetCreatedBy(user)
		}
	}
}

func assignUpdatedBy(db *gorm.DB) {
	if value, enabled := db.Statement.Model.(auditInterface); enabled {
		if user, ok := getCurrentUser(db); ok {
			value.SetUpdatedBy(user)
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) error {
	// register create callback
	createCallback := db.Callback().Create().Before("gorm:create")
	if err := createCallback.Register("audit:assign_created_by", assignCreatedBy); err != nil {
		return err
	}

	// register update callback
	updateCallback := db.Callback().Update().Before("gorm:update")
	if err := updateCallback.Register("audit:assign_updated_by", assignUpdatedBy); err != nil {
		return err
	}

	return nil
}
