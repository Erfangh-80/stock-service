package memory

import (
	"sync"

	pricehistorydomain "stock-service/internal/domain/price_history"
)

type PriceHistoryRepository struct {
	mu     sync.Mutex
	items  map[int64]*pricehistorydomain.PriceHistory
	nextID int64
}

func NewPriceHistoryRepository() *PriceHistoryRepository {
	return &PriceHistoryRepository{items: make(map[int64]*pricehistorydomain.PriceHistory)}
}

func (r *PriceHistoryRepository) Save(ph *pricehistorydomain.PriceHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ph.ID == 0 {
		r.nextID++
		ph.ID = r.nextID
	}
	r.items[ph.ID] = ph
	return nil
}

func (r *PriceHistoryRepository) FindByProductID(productID int32) ([]*pricehistorydomain.PriceHistory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*pricehistorydomain.PriceHistory
	for _, ph := range r.items {
		if ph.ProductID == productID {
			result = append(result, ph)
		}
	}
	return result, nil
}
