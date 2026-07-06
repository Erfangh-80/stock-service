package producttype

type Repository interface {
	Save(pt *ProductType) error
	FindByID(id int64) (*ProductType, error)
	FindByProductID(productID int32) ([]*ProductType, error)
	Delete(id int64) error
}
