package referenceprice

type ReferencePriceFilter struct {
	ProductID *int32
	Source    *string
	Page      int
	Limit     int
}

type Repository interface {
	Save(rp *ReferencePrice) error
	FindByID(id int64) (*ReferencePrice, error)
	FindByProductID(productID int32) (*ReferencePrice, error)
	FindAll(filter ReferencePriceFilter) ([]*ReferencePrice, int, error)
	Delete(id int64) error
}
