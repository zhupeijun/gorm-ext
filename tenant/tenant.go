package tenant

type Tenant struct {
	TenantID uint
}

func (model *Tenant) SetTenantID(tenantID uint) {
	model.TenantID = tenantID
}

func (model *Tenant) GetTenantID() uint {
	return model.TenantID
}
