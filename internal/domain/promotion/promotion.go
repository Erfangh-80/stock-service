package promotion

import "time"

// TODO: Promotion has no discount type/value field (%, fixed amount, BOGO, etc).
// Discount info currently lives only on ProductSale.final_price.
// Confirm if this entity should carry discount rules or remain metadata-only.
type PromotionStatus string

const (
	PromotionStatusInactive PromotionStatus = "inactive"
	PromotionStatusActive   PromotionStatus = "active"
)

type Promotion struct {
	ID                      int64
	Title                   string
	RequiresApproval        bool
	StartAt                 *time.Time
	EndAt                   *time.Time
	IsCountdown             bool
	ExpireSaleWithPromotion bool
	Status                  PromotionStatus
	CreatedAt               time.Time
}

func NewPromotion(title string) (*Promotion, error) {
	if err := ValidateTitle(title); err != nil {
		return nil, err
	}

	return &Promotion{
		Title:     title,
		Status:    PromotionStatusInactive,
		CreatedAt: time.Now(),
	}, nil
}

func (p *Promotion) Activate() {
	p.Status = PromotionStatusActive
}

func (p *Promotion) Deactivate() {
	p.Status = PromotionStatusInactive
}
