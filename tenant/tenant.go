package tenant

type Tenant struct {
	TenantID uint
}

func (model *Tenant) SetTenantID(tenantID uint) {
	model.TenantID = tenantID
}
