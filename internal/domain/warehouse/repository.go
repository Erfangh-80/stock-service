package warehouse

type WarehouseFilter struct {
	CreatedByUserID *int64
	IsPublic        *bool
	Page            int
	Limit           int
}

type Repository interface {
	Save(w *Warehouse) error
	FindByID(id int64) (*Warehouse, error)
	FindAll(filter WarehouseFilter) ([]*Warehouse, int, error)
	Delete(id int64) error
}
