package warehouse

import "time"

type Warehouse struct {
	ID             int64
	CreatedByUserID int64
	WarehouseName  string
	AddressID      *int64
	Phone          *string
	ContactPhone   *string
	IsPublic       bool
	CollectionMethod string
	CreatedAt      time.Time
}

func NewWarehouse(createdByUserID int64, warehouseName string) (*Warehouse, error) {
	if err := ValidateWarehouseName(warehouseName); err != nil {
		return nil, err
	}

	return &Warehouse{
		CreatedByUserID: createdByUserID,
		WarehouseName:   warehouseName,
		CreatedAt:       time.Now(),
	}, nil
}

func (w *Warehouse) MakePublic() {
	w.IsPublic = true
}

func (w *Warehouse) MakePrivate() {
	w.IsPublic = false
}

func (w *Warehouse) UpdatePhone(phone *string) {
	w.Phone = phone
}

func (w *Warehouse) UpdateContactPhone(phone *string) {
	w.ContactPhone = phone
}

func (w *Warehouse) UpdateCollectionMethod(method string) {
	w.CollectionMethod = method
}
