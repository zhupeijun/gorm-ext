package tenant

type TenantModel struct {
	TenantID *uint
}

func (model *TenantModel) SetTenantID(tenantID uint) {
	model.TenantID = &tenantID
}
