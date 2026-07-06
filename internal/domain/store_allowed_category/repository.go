package storeallowedcategory

type StoreCategoryFilter struct {
	StoreID *int64
	Page    int
	Limit   int
}

type Repository interface {
	Save(sac *StoreAllowedCategory) error
	FindByID(id int64) (*StoreAllowedCategory, error)
	FindAll(filter StoreCategoryFilter) ([]*StoreAllowedCategory, int, error)
	Delete(id int64) error
}
