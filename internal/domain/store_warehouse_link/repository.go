package storewarehouselink

type WarehouseLinkFilter struct {
	StoreID    *int64
	WarehouseID *int64
	Page       int
	Limit      int
}

type Repository interface {
	Save(swl *StoreWarehouseLink) error
	FindByID(id int64) (*StoreWarehouseLink, error)
	FindAll(filter WarehouseLinkFilter) ([]*StoreWarehouseLink, int, error)
	Delete(id int64) error
}
