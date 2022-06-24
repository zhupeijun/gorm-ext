package tenant

type TenantModel struct {
	TenantID *uint
}

func (model *TenantModel) SetTenantID(tenantID uint) {
	model.TenantID = &tenantID
}

// Deprecated: Will be replaced with TenantModel
type Tenant struct {
	TenantID uint
}

func (model *Tenant) SetTenantID(tenantID uint) {
	model.TenantID = tenantID
}
