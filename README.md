[![go report card](https://goreportcard.com/badge/github.com/zhupeijun/gorm-ext "go report card")](https://goreportcard.com/report/github.com/zhupeijun/gorm-ext)
[![Run Tests](https://github.com/zhupeijun/gorm-ext/actions/workflows/gorm-ext.yml/badge.svg)](https://github.com/zhupeijun/gorm-ext/actions/workflows/gorm-ext.yml)
[![GoDoc](https://godoc.org/github.com/zhupeijun/gorm-ext?status.svg)](https://godoc.org/github.com/zhupeijun/gorm-ext)

# Audit

Audit is used to record the last user who created/updated the model. It uses a `CreatedBy` and `UpdatedBy` field to save the information.  


# Tenant 

Tenant is used to record the model belong to which tenant. It uses `TenantID` to save the information when create the model.

# Example

### Register GORM callbacks

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
err = audit.RegisterCallbacks(db)
err = tenant.RegisterCallbacks(db)
```

### Making a model which include the extra information.

Embed `audit.Audit` for audit information and `tenant.Tenant` for tenant information. 

```go
type User struct {
    gorm.Model
    audit.Audit
    tenant.Tenant
    Name string
}
```

### Usage

```go
// setup the information to the db context
db = db.Set(tenant.ContextKeyTenantID, "tenant-1")
db = db.Set(audit.ContextKeyCurrentUser, "Admin")

// create a record
var user = User{Name: "Alice"}
db.Create(&user)
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).