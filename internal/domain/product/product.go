package product

import "time"

type ProductStatus string

const (
	ProductStatusPending  ProductStatus = "pending"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusRejected ProductStatus = "rejected"
	ProductStatusDeleted  ProductStatus = "deleted"
)

type OwnerType string

const (
	OwnerTypeSystem OwnerType = "system"
	OwnerTypeUser   OwnerType = "user"
)

type Product struct {
	ID               int32
	TitleFa          string
	TitleEn          *string
	Slug             string
	Description      *string
	BrandID          int64
	CategoryID       int64
	OwnerType        OwnerType
	OwnerID          *int64
	IsOriginal       bool
	Status           ProductStatus
	IsEnabled        bool
	EnabledAt        *time.Time
	DisabledAt       *time.Time
	MetaTitle        *string
	MetaDescription  *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	IndexImageFileID *int64
	DeletedAt        *time.Time
}

type ProductOption func(*Product)

func WithTitleEn(v *string) ProductOption {
	return func(p *Product) { p.TitleEn = v }
}

func WithDescription(v *string) ProductOption {
	return func(p *Product) { p.Description = v }
}

func WithOwnerType(v OwnerType) ProductOption {
	return func(p *Product) { p.OwnerType = v }
}

func WithOwnerID(v *int64) ProductOption {
	return func(p *Product) { p.OwnerID = v }
}

func WithIsOriginal(v bool) ProductOption {
	return func(p *Product) { p.IsOriginal = v }
}

func WithIndexImageFileID(v *int64) ProductOption {
	return func(p *Product) { p.IndexImageFileID = v }
}

func WithSlug(v string) ProductOption {
	return func(p *Product) { p.Slug = v }
}

func WithMetaTitle(v *string) ProductOption {
	return func(p *Product) { p.MetaTitle = v }
}

func WithMetaDescription(v *string) ProductOption {
	return func(p *Product) { p.MetaDescription = v }
}

func NewProduct(titleFa string, brandID, categoryID int64, opts ...ProductOption) (*Product, error) {
	if titleFa == "" {
		return nil, ErrTitleFaRequired
	}
	if brandID <= 0 {
		return nil, ErrInvalidBrandID
	}
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	now := time.Now()
	p := &Product{
		TitleFa:    titleFa,
		BrandID:    brandID,
		CategoryID: categoryID,
		OwnerType:  OwnerTypeSystem,
		IsOriginal: true,
		Status:     ProductStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p, nil
}

func (p *Product) MarkActive() {
	p.Status = ProductStatusActive
	p.Touch()
}

func (p *Product) MarkRejected() {
	p.Status = ProductStatusRejected
	p.Touch()
}

func (p *Product) SoftDelete() {
	p.Status = ProductStatusDeleted
	now := time.Now()
	p.DeletedAt = &now
	p.Touch()
}

func (p *Product) MarkEnabled() {
	p.IsEnabled = true
	now := time.Now()
	p.EnabledAt = &now
	p.DisabledAt = nil
	p.Touch()
}

func (p *Product) MarkDisabled() {
	p.IsEnabled = false
	now := time.Now()
	p.DisabledAt = &now
	p.EnabledAt = nil
	p.Touch()
}

func (p *Product) UpdateSEO(metaTitle, metaDescription *string) {
	p.MetaTitle = metaTitle
	p.MetaDescription = metaDescription
	p.Touch()
}

func (p *Product) GenerateSlug(slug string) {
	p.Slug = slug
	p.Touch()
}

func (p *Product) Touch() {
	p.UpdatedAt = time.Now()
}
