package tenant

type Tenant struct {
	TenantID string
}

func (model *Tenant) SetTenantID(tenantID string) {
	model.TenantID = tenantID
}

func (model *Tenant) GetTenantID() string {
	return model.TenantID
}
