package productinterface_test

import (
	"testing"
	"time"

	appproduct "stock-service/internal/application/product"
	"stock-service/internal/domain/product"
	iface "stock-service/internal/interface"
	productinterface "stock-service/internal/interface/product"
)

type mockCreateProduct struct{ fn func(appproduct.CreateProductInput) (*product.Product, error) }
func (m *mockCreateProduct) Execute(i appproduct.CreateProductInput) (*product.Product, error) { return m.fn(i) }

type mockGetProduct struct{ fn func(appproduct.GetProductInput) (*product.Product, error) }
func (m *mockGetProduct) Execute(i appproduct.GetProductInput) (*product.Product, error) { return m.fn(i) }

type mockUpdateProduct struct{ fn func(appproduct.UpdateProductInput) (*product.Product, error) }
func (m *mockUpdateProduct) Execute(i appproduct.UpdateProductInput) (*product.Product, error) { return m.fn(i) }

type mockActivateProduct struct{ fn func(appproduct.ActivateProductInput) (*product.Product, error) }
func (m *mockActivateProduct) Execute(i appproduct.ActivateProductInput) (*product.Product, error) { return m.fn(i) }

type mockRejectProduct struct{ fn func(appproduct.RejectProductInput) (*product.Product, error) }
func (m *mockRejectProduct) Execute(i appproduct.RejectProductInput) (*product.Product, error) { return m.fn(i) }

type mockSoftDeleteProduct struct{ fn func(appproduct.SoftDeleteProductInput) (*product.Product, error) }
func (m *mockSoftDeleteProduct) Execute(i appproduct.SoftDeleteProductInput) (*product.Product, error) { return m.fn(i) }

type mockEnableProduct struct{ fn func(appproduct.EnableProductInput) (*product.Product, error) }
func (m *mockEnableProduct) Execute(i appproduct.EnableProductInput) (*product.Product, error) { return m.fn(i) }

type mockDisableProduct struct{ fn func(appproduct.DisableProductInput) (*product.Product, error) }
func (m *mockDisableProduct) Execute(i appproduct.DisableProductInput) (*product.Product, error) { return m.fn(i) }

type mockUpdateSEO struct{ fn func(appproduct.UpdateSEOInput) (*product.Product, error) }
func (m *mockUpdateSEO) Execute(i appproduct.UpdateSEOInput) (*product.Product, error) { return m.fn(i) }

type mockListProducts struct{ fn func(appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error) }
func (m *mockListProducts) Execute(i appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error) { return m.fn(i) }

func nilMocks() (*mockEnableProduct, *mockDisableProduct, *mockUpdateSEO, *mockListProducts) {
	return &mockEnableProduct{nil}, &mockDisableProduct{nil}, &mockUpdateSEO{nil}, &mockListProducts{nil}
}

func baseProduct() *product.Product {
	return &product.Product{
		ID: 1, TitleFa: "محصول تست", BrandID: 1, CategoryID: 1,
		Status:    product.ProductStatusPending,
		OwnerType: product.OwnerTypeSystem,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{func(i appproduct.CreateProductInput) (*product.Product, error) { return p, nil }},
		&mockGetProduct{nil}, &mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	resp, err := adapter.Create(productinterface.CreateProductParams{
		TitleFa: "محصول تست", BrandID: 1, CategoryID: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.TitleFa != "محصول تست" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestAdapter_Create_ValidationError(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{func(i appproduct.CreateProductInput) (*product.Product, error) {
			return nil, product.ErrTitleFaRequired
		}},
		&mockGetProduct{nil}, &mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	_, err := adapter.Create(productinterface.CreateProductParams{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestAdapter_Get_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{func(i appproduct.GetProductInput) (*product.Product, error) {
			return p, nil
		}},
		&mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected id 1")
	}
}

func TestAdapter_Get_NotFound(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{func(i appproduct.GetProductInput) (*product.Product, error) {
			return nil, product.ErrProductNotFound
		}},
		&mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_Update_Success(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil},
		&mockUpdateProduct{func(i appproduct.UpdateProductInput) (*product.Product, error) {
			p := baseProduct()
			if i.TitleFa != nil {
				p.TitleFa = *i.TitleFa
			}
			return p, nil
		}},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	resp, err := adapter.Update(productinterface.UpdateProductParams{ID: 1, TitleFa: strPtr("updated")})
	if err != nil {
		t.Fatal(err)
	}
	if resp.TitleFa != "updated" {
		t.Error("expected updated title")
	}
}

func TestAdapter_Activate_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{func(i appproduct.ActivateProductInput) (*product.Product, error) { return p, nil }},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	resp, err := adapter.Activate(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected id 1")
	}
}

func TestAdapter_Reject_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil},
		&mockRejectProduct{func(i appproduct.RejectProductInput) (*product.Product, error) { return p, nil }},
		&mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	resp, err := adapter.Reject(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected id 1")
	}
}

func TestAdapter_SoftDelete_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil},
		&mockSoftDeleteProduct{func(i appproduct.SoftDeleteProductInput) (*product.Product, error) { return p, nil }},
		en, dis, seo, lst,
	)
	resp, err := adapter.SoftDelete(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected id 1")
	}
}

func TestAdapter_MapError_Default(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{func(i appproduct.CreateProductInput) (*product.Product, error) {
			return nil, product.ErrInvalidBrandID
		}},
		&mockGetProduct{nil}, &mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	_, err := adapter.Create(productinterface.CreateProductParams{})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func strPtr(s string) *string { return &s }
