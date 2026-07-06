package referenceprice_test

import (
	"testing"

	referencepriceapp "stock-service/internal/application/reference_price"
	"stock-service/internal/domain/reference_price"
)

func TestDeleteReferencePrice_Success(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewDeleteReferencePriceUseCase(repo)

	err := uc.Execute(referencepriceapp.DeleteReferencePriceInput{ID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rp, err := repo.FindByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if rp != nil {
		t.Error("expected reference price to be deleted")
	}
}

func TestDeleteReferencePrice_NotFound(t *testing.T) {
	repo := populateRepo()
	uc := referencepriceapp.NewDeleteReferencePriceUseCase(repo)

	err := uc.Execute(referencepriceapp.DeleteReferencePriceInput{ID: 999})
	if err != referenceprice.ErrReferencePriceNotFound {
		t.Errorf("expected ErrReferencePriceNotFound, got %v", err)
	}
}
