package productimage

import "time"

type ProductImage struct {
	ID        int64
	ProductID int32
	FileID    int64
	SortOrder int
	CreatedAt time.Time
}

func NewProductImage(productID int32, fileID int64, sortOrder int) (*ProductImage, error) {
	if productID <= 0 {
		return nil, ErrInvalidProductID
	}
	if fileID <= 0 {
		return nil, ErrInvalidFileID
	}

	return &ProductImage{
		ProductID: productID,
		FileID:    fileID,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}, nil
}

func (pi *ProductImage) UpdateSortOrder(order int) {
	pi.SortOrder = order
}
