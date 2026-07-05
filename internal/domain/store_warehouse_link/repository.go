package storewarehouselink

type Repository interface {
	Save(swl *StoreWarehouseLink) error
	FindByID(id int64) (*StoreWarehouseLink, error)
	Delete(id int64) error
}
