package productbundle

import "time"

type BundleType string

const (
	BundleRelation     BundleType = "bundle"
	UpsellRelation     BundleType = "upsell"
	CrossSellRelation  BundleType = "cross_sell"
)

type ProductBundle struct {
	ID               int64
	ProductID        int32
	RelatedProductID int32
	Type             BundleType
	SortOrder        int
	CreatedAt        time.Time
}

func NewProductBundle(productID, relatedProductID int32, bundleType BundleType, sortOrder int) (*ProductBundle, error) {
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if relatedProductID <= 0 {
		return nil, ErrInvalidRelatedProductID
	}
	if productID == relatedProductID {
		return nil, ErrSelfReference
	}
	if bundleType != BundleRelation && bundleType != UpsellRelation && bundleType != CrossSellRelation {
		return nil, ErrInvalidBundleType
	}

	return &ProductBundle{
		ProductID:        productID,
		RelatedProductID: relatedProductID,
		Type:             bundleType,
		SortOrder:        sortOrder,
		CreatedAt:        time.Now(),
	}, nil
}

func (pb *ProductBundle) UpdateSortOrder(order int) {
	pb.SortOrder = order
}
