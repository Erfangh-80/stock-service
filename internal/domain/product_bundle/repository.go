package productbundle

type Repository interface {
	Save(pb *ProductBundle) error
	FindByID(id int64) (*ProductBundle, error)
	FindByProductID(productID int32) ([]*ProductBundle, error)
	Delete(id int64) error
}
