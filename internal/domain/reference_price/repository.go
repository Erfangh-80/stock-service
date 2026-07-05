package referenceprice

type Repository interface {
	Save(rp *ReferencePrice) error
	FindByID(id int64) (*ReferencePrice, error)
	Delete(id int64) error
}
