package promotion

type PromotionFilter struct {
	Status       *PromotionStatus
	DiscountType *DiscountType
	Search       *string
	Page         int
	Limit        int
}

type Repository interface {
	Save(p *Promotion) error
	FindByID(id int64) (*Promotion, error)
	FindAll(filter PromotionFilter) ([]*Promotion, int, error)
	FindByCouponCode(code string) (*Promotion, error)
	Delete(id int64) error
}
