package inventory

type Repository interface {
	Save(inv *Inventory) error
	FindByID(id int64) (*Inventory, error)
	FindAll() ([]*Inventory, error)
	Delete(id int64) error
}
