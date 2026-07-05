package storeallowedcategory

type Repository interface {
	Save(sac *StoreAllowedCategory) error
	FindByID(id int64) (*StoreAllowedCategory, error)
	Delete(id int64) error
}
