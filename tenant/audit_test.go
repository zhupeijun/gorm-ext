package tenant

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Tenant
	Name string
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	err = RegisterCallbacks(db)
	if err != nil {
		panic(fmt.Sprintf("failed to register callback: %v", err))
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(fmt.Sprintf("failed auto migrate: %v", err))
	}
}

func TestAuditedCreate(t *testing.T) {
	db := db.Set(ContextKeyTenantID, "tenant-1")
	var user = User{Name: "Alice"}
	db.Create(&user)

	assert.Equal(t, "tenant-1", user.TenantID)
}
