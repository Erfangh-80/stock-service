package pricehistory

type Repository interface {
	Save(ph *PriceHistory) error
	FindByProductID(productID int32) ([]*PriceHistory, error)
}
