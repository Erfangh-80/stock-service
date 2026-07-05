package storeallowedcategory

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type StoreAllowedCategory struct {
	ID         int64
	StoreID    int64
	CategoryID int64
	Status     Status
	SupportNote string
	CreatedAt  time.Time
}

func NewStoreAllowedCategory(storeID, categoryID int64) *StoreAllowedCategory {
	return &StoreAllowedCategory{
		StoreID:    storeID,
		CategoryID: categoryID,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
	}
}

func (sac *StoreAllowedCategory) Approve() {
	sac.Status = StatusApproved
}

func (sac *StoreAllowedCategory) Reject() {
	sac.Status = StatusRejected
}
