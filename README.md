[![go report card](https://goreportcard.com/badge/github.com/zhupeijun/gorm-ext "go report card")](https://goreportcard.com/report/github.com/zhupeijun/gorm-ext)
[![Run Tests](https://github.com/zhupeijun/gorm-ext/actions/workflows/gorm-ext.yml/badge.svg)](https://github.com/zhupeijun/gorm-ext/actions/workflows/gorm-ext.yml)
[![GoDoc](https://godoc.org/github.com/zhupeijun/gorm-ext?status.svg)](https://godoc.org/github.com/zhupeijun/gorm-ext)

# Audit

Audit is used to record the last user who created/updated the model. It uses a `CreatedBy` and `UpdatedBy` field to save the information.  

## Example

### Register GORM callbacks

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
err = audit.RegisterCallbacks(db)
```

### Making a model which include the extra information.

Embed `audit.Audit` for audit information.

```go
type User struct {
    gorm.Model
    audit.Audit
    Name string
}
```

### Usage

```go
// setup the information to the db context
db = db.Set(audit.ContextKeyCurrentUser, "Admin")

// create a record
var user = User{Name: "Alice"}
db.Create(&user)
// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`created_by`,`updated_by`,`name`) VALUES ("...","...",NULL,"Admin","Admin","Alice")
```

# Tenant 

Tenant is used to record the model belong to which tenant. It uses `TenantID` to save the information when create the model.


## Example

### Register GORM callbacks

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
err = tenant.RegisterCallbacks(db)
```

### Making a model which include the extra information.

Embed `audit.Audit` for audit information.

```go
type User struct {
    gorm.Model
    tenant.Tenant
    Name string
    Role string
}
```

### Usage

```go
db := db.Set(ContextKeyTenantID, uint(1))

db.Create(&User{Name: "Alice", Role: "admin"})
// INSERT INTO `users` (`tenant_id`,`name`,`role`) VALUES (1,"Alice","admin") RETURNING `id`

db.Model(&User{}).Where("name", "Alice").Update("role", "user")
// UPDATE `users` SET `role`="user" WHERE `name` = "Alice" AND `tenant_id` = 1

var users []User
db.Find(&users)
// SELECT * FROM `users` WHERE `name` = "Alice" AND `tenant_id` = 1
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).