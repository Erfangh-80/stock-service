package productimage

type Repository interface {
	Save(image *ProductImage) error
	FindByID(id int64) (*ProductImage, error)
	FindByProductID(productID int32) ([]*ProductImage, error)
	Delete(id int64) error
}
