package inventory

import "time"

type SaleModel string

const SaleModelRetail SaleModel = "retail"

type PromotionStatus string

const PromotionStatusPending PromotionStatus = "pending"

type Condition string

const ConditionNew Condition = "new"

type VendorSaleStatus string

const (
	VendorSaleStatusActive    VendorSaleStatus = "active"
	VendorSaleStatusSuspended VendorSaleStatus = "suspended"
	VendorSaleStatusClosed    VendorSaleStatus = "closed"
)

type SystemSaleStatus string

const (
	SystemSaleStatusActive    SystemSaleStatus = "active"
	SystemSaleStatusSuspended SystemSaleStatus = "suspended"
	SystemSaleStatusClosed    SystemSaleStatus = "closed"
)

type Inventory struct {
	ID               int64
	StoreID          int64
	WarehouseID      int64
	ProductID        int32
	SaleModel        SaleModel
	BasePrice        float64
	PromotionID      *int64
	FinalPrice       *float64
	StartAt          *time.Time
	EndAt            *time.Time
	PromotionStatus  PromotionStatus
	Attributes       map[string]any
	InstantQty       int
	ScheduledQty     map[string]int
	MinOrderQty      int
	MaxOrderQty      *int
	Condition        Condition
	VendorSaleStatus VendorSaleStatus
	SystemSaleStatus SystemSaleStatus
	CreatedAt        time.Time
}

func NewInventory(storeID, warehouseID int64, productID int32, basePrice float64) (*Inventory, error) {
	if err := ValidateBasePrice(basePrice); err != nil {
		return nil, err
	}

	return &Inventory{
		StoreID:          storeID,
		WarehouseID:      warehouseID,
		ProductID:        productID,
		SaleModel:        SaleModelRetail,
		BasePrice:        basePrice,
		PromotionStatus:  PromotionStatusPending,
		InstantQty:       0,
		MinOrderQty:      1,
		Condition:        ConditionNew,
		VendorSaleStatus: VendorSaleStatusActive,
		SystemSaleStatus: SystemSaleStatusActive,
		CreatedAt:        time.Now(),
	}, nil
}

func (inv *Inventory) ApplyPromotion(promotionID int64, finalPrice float64, startAt, endAt time.Time) error {
	if inv.PromotionID != nil {
		return ErrPromotionAlreadyApplied
	}

	if err := ValidatePromotionDates(startAt, endAt); err != nil {
		return err
	}

	if err := ValidateFinalPrice(finalPrice); err != nil {
		return err
	}

	inv.PromotionID = &promotionID
	inv.FinalPrice = &finalPrice
	inv.StartAt = &startAt
	inv.EndAt = &endAt
	inv.PromotionStatus = PromotionStatusPending

	return nil
}

func (inv *Inventory) RemovePromotion() error {
	if inv.PromotionID == nil {
		return ErrNoActivePromotion
	}

	inv.PromotionID = nil
	inv.FinalPrice = nil
	inv.StartAt = nil
	inv.EndAt = nil
	inv.PromotionStatus = PromotionStatusPending

	return nil
}

func (inv *Inventory) UpdateInventory(instantQty int, scheduledQty map[string]int, minOrderQty int, maxOrderQty *int) error {
	if err := ValidateInstantQty(instantQty); err != nil {
		return err
	}

	if err := ValidateMinOrderQty(minOrderQty); err != nil {
		return err
	}

	if maxOrderQty != nil {
		if err := ValidateMaxOrderQty(minOrderQty, *maxOrderQty); err != nil {
			return err
		}
	}

	inv.InstantQty = instantQty
	inv.ScheduledQty = scheduledQty
	inv.MinOrderQty = minOrderQty
	inv.MaxOrderQty = maxOrderQty

	return nil
}

func (inv *Inventory) SuspendVendorSale() error {
	if inv.VendorSaleStatus == VendorSaleStatusClosed {
		return ErrVendorSaleStatusTransition
	}
	inv.VendorSaleStatus = VendorSaleStatusSuspended
	return nil
}

func (inv *Inventory) CloseVendorSale() error {
	if inv.VendorSaleStatus == VendorSaleStatusClosed {
		return ErrVendorSaleStatusTransition
	}
	inv.VendorSaleStatus = VendorSaleStatusClosed
	return nil
}

func (inv *Inventory) SuspendSystemSale() error {
	if inv.SystemSaleStatus == SystemSaleStatusClosed {
		return ErrSystemSaleStatusTransition
	}
	inv.SystemSaleStatus = SystemSaleStatusSuspended
	return nil
}

func (inv *Inventory) CloseSystemSale() error {
	if inv.SystemSaleStatus == SystemSaleStatusClosed {
		return ErrSystemSaleStatusTransition
	}
	inv.SystemSaleStatus = SystemSaleStatusClosed
	return nil
}

func (inv *Inventory) ReserveQuantity(qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	if inv.InstantQty < qty {
		return ErrInsufficientStock
	}
	inv.InstantQty -= qty
	return nil
}

func (inv *Inventory) ReleaseQuantity(qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	inv.InstantQty += qty
	return nil
}

func (inv *Inventory) HasLowStock(threshold int) bool {
	return inv.InstantQty <= threshold
}

func (inv *Inventory) ValidateScheduledQty(deliveryDate string, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	if err := ValidateScheduledDeliveryDate(deliveryDate); err != nil {
		return err
	}
	return nil
}
