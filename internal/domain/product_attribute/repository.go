package productattribute

type Repository interface {
	Save(pa *ProductAttribute) error
	FindByID(id int64) (*ProductAttribute, error)
	FindByProductID(productID int32) ([]*ProductAttribute, error)
	Delete(id int64) error
}
