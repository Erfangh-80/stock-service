package warehouse

type Repository interface {
	Save(w *Warehouse) error
	FindByID(id int64) (*Warehouse, error)
	Delete(id int64) error
}
