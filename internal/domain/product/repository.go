package product

type Repository interface {
	FindByID(id int32) (*Product, error)
	FindByTitle(query string) ([]*Product, error)
	Save(product *Product) error
}
