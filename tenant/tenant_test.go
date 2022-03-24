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
	Role string
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

func createUser() {
	db.Unscoped().Where("1=1").Delete(&User{})

	db.Set(ContextKeyTenantID, "tenant-1").Create(&User{Name: "Alice", Role: "admin"})
	db.Set(ContextKeyTenantID, "tenant-2").Create(&User{Name: "Alice", Role: "admin"})

}

func TestUserCreateAndRead(t *testing.T) {
	createUser()

	tenant1DB := db.Set(ContextKeyTenantID, "tenant-1")
	tenant2DB := db.Set(ContextKeyTenantID, "tenant-2")

	var count int64

	tenant1DB.Model(&User{}).Where("tenant_id", "tenant-1").Count(&count)
	assert.Equal(t, int64(1), count)

	tenant2DB.Model(&User{}).Where("tenant_id", "tenant-2").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestUserUpdate(t *testing.T) {
	createUser()

	var user1 User
	tenant1DB := db.Set(ContextKeyTenantID, "tenant-1")
	tenant2DB := db.Set(ContextKeyTenantID, "tenant-2")

	tenant1DB.Model(&User{}).Where("name", "Alice").Update("role", "user")
	tenant1DB.Where("name", "Alice").Find(&user1)
	assert.Equal(t, "user", user1.Role)

	var user2 User
	tenant2DB.Where("name", "Alice").Find(&user2)
	assert.Equal(t, "admin", user2.Role)
}

func TestUserDelete(t *testing.T) {
	createUser()

	tenant1DB := db.Set(ContextKeyTenantID, "tenant-1")
	tenant2DB := db.Set(ContextKeyTenantID, "tenant-2")

	tenant1DB.Where("name", "Alice").Delete(&User{})

	var count int64
	tenant1DB.Model(&User{}).Where("name", "Alice").Count(&count)
	assert.Equal(t, int64(0), count)

	tenant2DB.Model(&User{}).Where("name", "Alice").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestUserWithoutScope(t *testing.T) {
	createUser()

	var count int64
	db.Set(DisableTenantScope, true).Model(&User{}).Count(&count)
	assert.Equal(t, int64(2), count)
}
