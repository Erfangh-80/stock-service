package salescommission

type SalesCommissionFilter struct {
	InventoryID *int64
	SaleModel   *string
	Page        int
	Limit       int
}

type Repository interface {
	Save(sc *SalesCommission) error
	FindByID(id int64) (*SalesCommission, error)
	FindByInventoryID(inventoryID int64) (*SalesCommission, error)
	FindAll(filter SalesCommissionFilter) ([]*SalesCommission, int, error)
	Delete(id int64) error
}
