package store

import "time"

type StoreStatus string
const StoreStatusActive StoreStatus = "active"

type Store struct {
	ID                     int64
	UserID                 int64
	StoreName              string
	Status                 StoreStatus
	AddressID              *int64
	ContactPhone           *string
	MediaAssets            map[string]any
	IsCommissionApplicable bool
	IsBulkSaleEnabled      bool
	CreatedAt              time.Time
}

func NewStore(userID int64, storeName string) (*Store, error) {
	if err := ValidateStoreName(storeName); err != nil {
		return nil, err
	}

	return &Store{
		UserID:                 userID,
		StoreName:              storeName,
		Status:                 StoreStatusActive,
		IsCommissionApplicable: true,
		IsBulkSaleEnabled:      false,
		CreatedAt:              time.Now(),
	}, nil
}

func (s *Store) EnableBulkSale() {
	s.IsBulkSaleEnabled = true
}

func (s *Store) DisableBulkSale() {
	s.IsBulkSaleEnabled = false
}

func (s *Store) EnableCommission() {
	s.IsCommissionApplicable = true
}

func (s *Store) DisableCommission() {
	s.IsCommissionApplicable = false
}

func (s *Store) UpdateName(name string) error {
	if err := ValidateStoreName(name); err != nil {
		return err
	}
	s.StoreName = name
	return nil
}

func (s *Store) UpdateContactInfo(phone *string) {
	s.ContactPhone = phone
}
