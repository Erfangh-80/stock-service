package referenceprice_test

import (
	"testing"

	referencepriceapp "stock-service/internal/application/reference_price"
	"stock-service/internal/domain/reference_price"
)

func populateRepo() *inMemoryReferencePriceRepo {
	repo := newInMemoryReferencePriceRepo()
	repo.Save(&referenceprice.ReferencePrice{
		ProductID: 1, Price: 99.99, Source: "supplier",
	})
	return repo
}

func TestGetReferencePrice_Success(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewGetReferencePriceUseCase(repo)

	rp, err := uc.Execute(referencepriceapp.GetReferencePriceInput{ID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rp.ProductID != 1 || rp.Price != 99.99 || rp.Source != "supplier" {
		t.Errorf("unexpected reference price: %+v", rp)
	}
}

func TestGetReferencePrice_NotFound(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewGetReferencePriceUseCase(repo)

	_, err := uc.Execute(referencepriceapp.GetReferencePriceInput{ID: 999})
	if err != referenceprice.ErrReferencePriceNotFound {
		t.Errorf("expected ErrReferencePriceNotFound, got %v", err)
	}
}

func TestGetByProductReferencePrice_Success(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewGetByProductReferencePriceUseCase(repo)

	rp, err := uc.Execute(referencepriceapp.GetByProductReferencePriceInput{ProductID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rp.ProductID != 1 || rp.Price != 99.99 || rp.Source != "supplier" {
		t.Errorf("unexpected reference price: %+v", rp)
	}
}

func TestGetByProductReferencePrice_NotFound(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewGetByProductReferencePriceUseCase(repo)

	_, err := uc.Execute(referencepriceapp.GetByProductReferencePriceInput{ProductID: 999})
	if err != referenceprice.ErrReferencePriceNotFound {
		t.Errorf("expected ErrReferencePriceNotFound, got %v", err)
	}
}
