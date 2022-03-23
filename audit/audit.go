package audit

type Audit struct {
	CreatedBy string
	UpdatedBy string
}

func (model *Audit) SetCreatedBy(createdBy string) {
	model.CreatedBy = createdBy
}

func (model *Audit) GetCreatedBy() string {
	return model.GetCreatedBy()
}

func (model *Audit) SetUpdatedBy(updatedBy string) {
	model.UpdatedBy = updatedBy
}

func (model *Audit) GetUpdatedBy() string {
	return model.GetUpdatedBy()
}
