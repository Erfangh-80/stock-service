package product_test

import (
	"testing"

	"stock-service/internal/domain/product"
)

func TestNewProduct_Success(t *testing.T) {
	t.Parallel()

	titleFa := "محصول آزمایشی"
	brandID := int64(10)
	categoryID := int64(20)

	p, err := product.NewProduct(titleFa, brandID, categoryID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.TitleFa != titleFa {
		t.Errorf("expected TitleFa %q, got %q", titleFa, p.TitleFa)
	}
	if p.BrandID != brandID {
		t.Errorf("expected BrandID %d, got %d", brandID, p.BrandID)
	}
	if p.CategoryID != categoryID {
		t.Errorf("expected CategoryID %d, got %d", categoryID, p.CategoryID)
	}
	if p.TitleEn != nil {
		t.Error("expected TitleEn to be nil")
	}
	if p.Description != nil {
		t.Error("expected Description to be nil")
	}
	if p.OwnerType != product.OwnerTypeSystem {
		t.Errorf("expected OwnerType %q, got %q", product.OwnerTypeSystem, p.OwnerType)
	}
	if p.OwnerID != nil {
		t.Error("expected OwnerID to be nil")
	}
	if !p.IsOriginal {
		t.Error("expected IsOriginal to be true")
	}
	if p.Status != product.ProductStatusPending {
		t.Errorf("expected Status %q, got %q", product.ProductStatusPending, p.Status)
	}
	if p.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if p.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
	if p.IndexImageFileID != nil {
		t.Error("expected IndexImageFileID to be nil")
	}
}

func TestNewProduct_WithOptionalFields(t *testing.T) {
	t.Parallel()

	titleEn := "test product"
	desc := "a description"
	ownerID := int64(5)
	imageID := int64(99)

	p, err := product.NewProduct("test", 1, 2,
		product.WithTitleEn(&titleEn),
		product.WithDescription(&desc),
		product.WithOwnerType(product.OwnerTypeUser),
		product.WithOwnerID(&ownerID),
		product.WithIsOriginal(false),
		product.WithIndexImageFileID(&imageID),
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.TitleEn == nil || *p.TitleEn != titleEn {
		t.Errorf("expected TitleEn %q, got %v", titleEn, p.TitleEn)
	}
	if p.Description == nil || *p.Description != desc {
		t.Errorf("expected Description %q, got %v", desc, p.Description)
	}
	if p.OwnerType != product.OwnerTypeUser {
		t.Errorf("expected OwnerType %q, got %q", product.OwnerTypeUser, p.OwnerType)
	}
	if p.OwnerID == nil || *p.OwnerID != ownerID {
		t.Errorf("expected OwnerID %d, got %v", ownerID, p.OwnerID)
	}
	if p.IsOriginal {
		t.Error("expected IsOriginal to be false")
	}
	if p.IndexImageFileID == nil || *p.IndexImageFileID != imageID {
		t.Errorf("expected IndexImageFileID %d, got %v", imageID, p.IndexImageFileID)
	}
}

func TestNewProduct_EmptyTitleFa_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := product.NewProduct("", 1, 2)
	if err != product.ErrTitleFaRequired {
		t.Errorf("expected ErrTitleFaRequired, got %v", err)
	}
}

func TestNewProduct_ZeroBrandID_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := product.NewProduct("test", 0, 2)
	if err != product.ErrInvalidBrandID {
		t.Errorf("expected ErrInvalidBrandID, got %v", err)
	}
}

func TestNewProduct_ZeroCategoryID_ReturnsError(t *testing.T) {
	t.Parallel()

	_, err := product.NewProduct("test", 1, 0)
	if err != product.ErrInvalidCategoryID {
		t.Errorf("expected ErrInvalidCategoryID, got %v", err)
	}
}

func TestSoftDelete_Success(t *testing.T) {
	t.Parallel()

	p, err := product.NewProduct("test", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	p.SoftDelete()

	if p.Status != product.ProductStatusDeleted {
		t.Errorf("expected Status %q, got %q", product.ProductStatusDeleted, p.Status)
	}
	if p.DeletedAt == nil || p.DeletedAt.IsZero() {
		t.Error("expected DeletedAt to be set")
	}
}

func TestMarkActive_Success(t *testing.T) {
	t.Parallel()

	p, err := product.NewProduct("test", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	p.MarkActive()

	if p.Status != product.ProductStatusActive {
		t.Errorf("expected Status %q, got %q", product.ProductStatusActive, p.Status)
	}
}

func TestActivate(t *testing.T) {
	t.Parallel()

	p, err := product.NewProduct("test", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	p.MarkActive()
	if p.Status != product.ProductStatusActive {
		t.Errorf("expected Status %q, got %q", product.ProductStatusActive, p.Status)
	}

	p.MarkRejected()
	if p.Status != product.ProductStatusRejected {
		t.Errorf("expected Status %q, got %q", product.ProductStatusRejected, p.Status)
	}
}

func TestUpdateTimestamps(t *testing.T) {
	t.Parallel()

	p, err := product.NewProduct("test", 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	old := p.UpdatedAt
	p.Touch()

	if !p.UpdatedAt.After(old) {
		t.Error("expected UpdatedAt to be after the old value")
	}
}
