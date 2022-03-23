package audit

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
	Audit
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
	db := db.Set(ContextKeyCurrentUser, "Admin")
	var user = User{Name: "Alice"}
	db.Create(&user)

	assert.Equal(t, "Admin", user.CreatedBy)
	assert.Equal(t, "Admin", user.UpdatedBy)
}

func TestAuditedUpdate(t *testing.T) {
	db := db.Set(ContextKeyCurrentUser, "Admin")
	var user = User{Name: "Alice"}
	db.Create(&user)

	db = db.Set(ContextKeyCurrentUser, "Operator")
	user.Name = "Bob"
	db.Model(&user).Update("name", "Bob")

	assert.Equal(t, "Bob", user.Name)
	assert.Equal(t, "Admin", user.CreatedBy)
	assert.Equal(t, "Operator", user.UpdatedBy)
}
