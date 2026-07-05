package product_test

import (
	"strings"
	"testing"

	"stock-service/internal/domain/product"
	appproduct "stock-service/internal/application/product"
)

type inMemoryProductRepo struct {
	products map[int32]*product.Product
	nextID   int32
}

func newInMemoryProductRepo() *inMemoryProductRepo {
	return &inMemoryProductRepo{
		products: make(map[int32]*product.Product),
		nextID:   1,
	}
}

func (r *inMemoryProductRepo) Save(p *product.Product) error {
	if p.ID == 0 {
		p.ID = r.nextID
		r.nextID++
	}
	r.products[p.ID] = p
	return nil
}

func (r *inMemoryProductRepo) FindByID(id int32) (*product.Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (r *inMemoryProductRepo) FindByTitle(query string) ([]*product.Product, error) {
	var result []*product.Product
	for _, p := range r.products {
		if strings.Contains(p.TitleFa, query) {
			result = append(result, p)
		}
	}
	return result, nil
}

func TestCreateProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	uc := appproduct.NewCreateProductUseCase(repo)

	input := appproduct.CreateProductInput{
		TitleFa:    "محصول آزمایشی",
		BrandID:    10,
		CategoryID: 20,
	}

	p, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID == 0 {
		t.Error("expected ID to be set")
	}
	if p.TitleFa != "محصول آزمایشی" {
		t.Errorf("expected TitleFa %q, got %q", "محصول آزمایشی", p.TitleFa)
	}
	if p.BrandID != 10 {
		t.Errorf("expected BrandID 10, got %d", p.BrandID)
	}
	if p.CategoryID != 20 {
		t.Errorf("expected CategoryID 20, got %d", p.CategoryID)
	}
	if p.Status != product.ProductStatusPending {
		t.Errorf("expected Status %q, got %q", product.ProductStatusPending, p.Status)
	}
}

func TestCreateProduct_WithOptionals(t *testing.T) {
	repo := newInMemoryProductRepo()
	uc := appproduct.NewCreateProductUseCase(repo)

	titleEn := "test product"
	desc := "a description"
	ownerID := int64(5)
	imageID := int64(99)
	isOriginal := false

	input := appproduct.CreateProductInput{
		TitleFa:          "test",
		BrandID:          1,
		CategoryID:       2,
		TitleEn:          &titleEn,
		Description:      &desc,
		OwnerType:        product.OwnerTypeUser,
		OwnerID:          &ownerID,
		IsOriginal:       &isOriginal,
		IndexImageFileID: &imageID,
	}

	p, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.TitleEn == nil || *p.TitleEn != titleEn {
		t.Errorf("expected TitleEn %q, got %v", titleEn, p.TitleEn)
	}
	if p.OwnerType != product.OwnerTypeUser {
		t.Errorf("expected OwnerType %q, got %q", product.OwnerTypeUser, p.OwnerType)
	}
	if p.IsOriginal {
		t.Error("expected IsOriginal to be false")
	}
}

func TestCreateProduct_EmptyTitleFa_ReturnsError(t *testing.T) {
	repo := newInMemoryProductRepo()
	uc := appproduct.NewCreateProductUseCase(repo)

	_, err := uc.Execute(appproduct.CreateProductInput{
		TitleFa:    "",
		BrandID:    1,
		CategoryID: 2,
	})
	if err != product.ErrTitleFaRequired {
		t.Errorf("expected ErrTitleFaRequired, got %v", err)
	}
}

func TestGetProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	createUC := appproduct.NewCreateProductUseCase(repo)
	getUC := appproduct.NewGetProductUseCase(repo)

	created, err := createUC.Execute(appproduct.CreateProductInput{
		TitleFa: "test", BrandID: 1, CategoryID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := getUC.Execute(appproduct.GetProductInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	repo := newInMemoryProductRepo()
	getUC := appproduct.NewGetProductUseCase(repo)

	_, err := getUC.Execute(appproduct.GetProductInput{ID: 999})
	if err != product.ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	createUC := appproduct.NewCreateProductUseCase(repo)
	updateUC := appproduct.NewUpdateProductUseCase(repo)

	created, err := createUC.Execute(appproduct.CreateProductInput{
		TitleFa: "old title", BrandID: 1, CategoryID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	titleEn := "new english"
	desc := "new desc"
	newTitleFa := "عنوان جدید"

	updated, err := updateUC.Execute(appproduct.UpdateProductInput{
		ID:          created.ID,
		TitleFa:     &newTitleFa,
		TitleEn:     &titleEn,
		Description: &desc,
		BrandID:     ptr(int64(5)),
		CategoryID:  ptr(int64(6)),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.TitleFa != newTitleFa {
		t.Errorf("expected TitleFa %q, got %q", newTitleFa, updated.TitleFa)
	}
	if updated.TitleEn == nil || *updated.TitleEn != titleEn {
		t.Errorf("expected TitleEn %q, got %v", titleEn, updated.TitleEn)
	}
	if updated.BrandID != 5 {
		t.Errorf("expected BrandID 5, got %d", updated.BrandID)
	}
	if updated.CategoryID != 6 {
		t.Errorf("expected CategoryID 6, got %d", updated.CategoryID)
	}
}

func TestUpdateProduct_NotFound(t *testing.T) {
	repo := newInMemoryProductRepo()
	updateUC := appproduct.NewUpdateProductUseCase(repo)

	_, err := updateUC.Execute(appproduct.UpdateProductInput{ID: 999})
	if err != product.ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestActivateProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	createUC := appproduct.NewCreateProductUseCase(repo)
	activateUC := appproduct.NewActivateProductUseCase(repo)

	created, err := createUC.Execute(appproduct.CreateProductInput{
		TitleFa: "test", BrandID: 1, CategoryID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	p, err := activateUC.Execute(appproduct.ActivateProductInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Status != product.ProductStatusActive {
		t.Errorf("expected Status %q, got %q", product.ProductStatusActive, p.Status)
	}
}

func TestRejectProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	createUC := appproduct.NewCreateProductUseCase(repo)
	rejectUC := appproduct.NewRejectProductUseCase(repo)

	created, err := createUC.Execute(appproduct.CreateProductInput{
		TitleFa: "test", BrandID: 1, CategoryID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	p, err := rejectUC.Execute(appproduct.RejectProductInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Status != product.ProductStatusRejected {
		t.Errorf("expected Status %q, got %q", product.ProductStatusRejected, p.Status)
	}
}

func TestSoftDeleteProduct_Success(t *testing.T) {
	repo := newInMemoryProductRepo()
	createUC := appproduct.NewCreateProductUseCase(repo)
	deleteUC := appproduct.NewSoftDeleteProductUseCase(repo)

	created, err := createUC.Execute(appproduct.CreateProductInput{
		TitleFa: "test", BrandID: 1, CategoryID: 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	p, err := deleteUC.Execute(appproduct.SoftDeleteProductInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Status != product.ProductStatusDeleted {
		t.Errorf("expected Status %q, got %q", product.ProductStatusDeleted, p.Status)
	}
	if p.DeletedAt == nil || p.DeletedAt.IsZero() {
		t.Error("expected DeletedAt to be set")
	}
}

func ptr[T any](v T) *T { return &v }
