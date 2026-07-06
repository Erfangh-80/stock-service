package product

type ProductFilter struct {
	OwnerType  *OwnerType
	OwnerID    *int64
	Status     *ProductStatus
	CategoryID *int64
	BrandID    *int64
	Search     *string
	Page       int
	Limit      int
}

type Repository interface {
	FindByID(id int32) (*Product, error)
	FindByTitle(query string) ([]*Product, error)
	FindAll(filter ProductFilter) ([]*Product, error)
	Count(filter ProductFilter) (int, error)
	Save(product *Product) error
}

