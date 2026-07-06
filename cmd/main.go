package main

import (
	"log"
	"net/http"

	appbrand "stock-service/internal/application/brand"
	appcategory "stock-service/internal/application/category"
	appinventory "stock-service/internal/application/inventory"
	appproduct "stock-service/internal/application/product"
	appproductimage "stock-service/internal/application/product_image"
	appproducttype "stock-service/internal/application/product_type"
	appproductattribute "stock-service/internal/application/product_attribute"
	apppricehistory "stock-service/internal/application/price_history"
	appproductbundle "stock-service/internal/application/product_bundle"
	apppromotion "stock-service/internal/application/promotion"
	appreferenceprice "stock-service/internal/application/reference_price"
	appsalescommission "stock-service/internal/application/sales_commission"
	appstore "stock-service/internal/application/store"
	appstoreallowedcategory "stock-service/internal/application/store_allowed_category"
	appstorewarehouselink "stock-service/internal/application/store_warehouse_link"
	appwarehouse "stock-service/internal/application/warehouse"

	"stock-service/internal/infrastructure/memory"

	warehousedomain "stock-service/internal/domain/warehouse"
	storewarehouselinkdomain "stock-service/internal/domain/store_warehouse_link"

	brandinterface "stock-service/internal/interface/brand"
	categoryinterface "stock-service/internal/interface/category"
	inventoryinterface "stock-service/internal/interface/inventory"
	productinterface "stock-service/internal/interface/product"
	productimageinterface "stock-service/internal/interface/product_image"
	producttypeinterface "stock-service/internal/interface/product_type"
	productattributeinterface "stock-service/internal/interface/product_attribute"
	pricehistoryinterface "stock-service/internal/interface/price_history"
	productbundleinterface "stock-service/internal/interface/product_bundle"
	promotioninterface "stock-service/internal/interface/promotion"
	referencepriceinterface "stock-service/internal/interface/reference_price"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
	storeinterface "stock-service/internal/interface/store"
	storeallowedcategoryinterface "stock-service/internal/interface/store_allowed_category"
	storewarehouselinkinterface "stock-service/internal/interface/store_warehouse_link"
	warehouseinterface "stock-service/internal/interface/warehouse"

	"stock-service/internal/interface/http/router"
)

func main() {
	brandRepo := memory.NewBrandRepository()
	categoryRepo := memory.NewCategoryRepository()
	storeRepo := memory.NewStoreRepository()
	inventoryRepo := memory.NewInventoryRepository()
	promotionRepo := memory.NewPromotionRepository()
	warehouseRepo := memory.NewWarehouseRepository()
	refPriceRepo := memory.NewReferencePriceRepository()
	salesCommRepo := memory.NewSalesCommissionRepository()
	catCommRuleRepo := memory.NewCategoryCommissionRuleRepository()
	storeCatRepo := memory.NewStoreCategoryRepository()
	warehouseLinkRepo := memory.NewWarehouseLinkRepository()
	productRepo := memory.NewProductRepository()
	memory.SeedProducts(productRepo)
	productImageRepo := memory.NewProductImageRepository()
	productTypeRepo := memory.NewProductTypeRepository()
	productAttrRepo := memory.NewProductAttributeRepository()
	priceHistRepo := memory.NewPriceHistoryRepository()
	productBundleRepo := memory.NewProductBundleRepository()

	brandAdapter := brandinterface.NewAdapter(
		*appbrand.NewCreateBrandUseCase(brandRepo),
		*appbrand.NewGetBrandUseCase(brandRepo),
		*appbrand.NewUpdateBrandUseCase(brandRepo),
		*appbrand.NewDeleteBrandUseCase(brandRepo),
		*appbrand.NewListBrandsUseCase(brandRepo),
	)

	categoryAdapter := categoryinterface.NewAdapter(
		*appcategory.NewCreateCategoryUseCase(categoryRepo),
		*appcategory.NewGetCategoryUseCase(categoryRepo),
		*appcategory.NewUpdateCategoryUseCase(categoryRepo),
		*appcategory.NewDeleteCategoryUseCase(categoryRepo),
		*appcategory.NewListCategoriesUseCase(categoryRepo),
	)

	storeAdapter := storeinterface.NewAdapter(
		appstore.NewCreateStoreUseCase(storeRepo),
		appstore.NewGetStoreUseCase(storeRepo),
		appstore.NewListStoresUseCase(storeRepo),
		appstore.NewToggleBulkSaleUseCase(storeRepo),
		appstore.NewToggleCommissionUseCase(storeRepo),
		appstore.NewUpdateContactUseCase(storeRepo),
		appstore.NewUpdateStoreNameUseCase(storeRepo),
		appstore.NewUpdateStoreProfileUseCase(storeRepo),
		appstore.NewDeleteStoreUseCase(storeRepo),
	)

	inventoryAdapter := inventoryinterface.NewAdapter(
		appinventory.NewCreateInventoryUseCase(inventoryRepo, productRepo),
		appinventory.NewGetInventoryUseCase(inventoryRepo),
		appinventory.NewListInventoryUseCase(inventoryRepo),
		appinventory.NewDeleteInventoryUseCase(inventoryRepo),
		appinventory.NewSearchInventoryUseCase(inventoryRepo, productRepo),
		appinventory.NewApplyPromotionUseCase(inventoryRepo, promotionRepo),
		appinventory.NewRemovePromotionUseCase(inventoryRepo),
		appinventory.NewUpdateInventoryUseCase(inventoryRepo),
		appinventory.NewSuspendVendorSaleUseCase(inventoryRepo),
		appinventory.NewCloseVendorSaleUseCase(inventoryRepo),
		appinventory.NewSuspendSystemSaleUseCase(inventoryRepo),
		appinventory.NewCloseSystemSaleUseCase(inventoryRepo),
		appinventory.NewReserveQuantityUseCase(inventoryRepo),
		appinventory.NewReleaseQuantityUseCase(inventoryRepo),
		appinventory.NewCheckLowStockUseCase(inventoryRepo),
	)

	promotionAdapter := promotioninterface.NewAdapter(
		apppromotion.NewCreatePromotionUseCase(promotionRepo),
		apppromotion.NewGetPromotionUseCase(promotionRepo),
		apppromotion.NewActivatePromotionUseCase(promotionRepo),
		apppromotion.NewDeactivatePromotionUseCase(promotionRepo),
		apppromotion.NewUpdatePromotionUseCase(promotionRepo),
		apppromotion.NewDeletePromotionUseCase(promotionRepo),
		apppromotion.NewListPromotionsUseCase(promotionRepo),
	)

	productAdapter := productinterface.NewAdapter(
		appproduct.NewCreateProductUseCase(productRepo),
		appproduct.NewGetProductUseCase(productRepo),
		appproduct.NewUpdateProductUseCase(productRepo),
		appproduct.NewActivateProductUseCase(productRepo),
		appproduct.NewRejectProductUseCase(productRepo),
		appproduct.NewSoftDeleteProductUseCase(productRepo),
		appproduct.NewEnableProductUseCase(productRepo),
		appproduct.NewDisableProductUseCase(productRepo),
		appproduct.NewUpdateSEOUseCase(productRepo),
		appproduct.NewListProductsUseCase(productRepo),
	)

	productImageAdapter := productimageinterface.NewAdapter(
		*appproductimage.NewCreateImageUseCase(productImageRepo),
		*appproductimage.NewListImagesUseCase(productImageRepo),
		*appproductimage.NewDeleteImageUseCase(productImageRepo),
	)

	productTypeAdapter := producttypeinterface.NewAdapter(
		*appproducttype.NewCreateTypeUseCase(productTypeRepo),
		*appproducttype.NewListTypesUseCase(productTypeRepo),
	)

	productAttrAdapter := productattributeinterface.NewAdapter(
		*appproductattribute.NewCreateAttributeUseCase(productAttrRepo),
		*appproductattribute.NewListAttributesUseCase(productAttrRepo),
	)

	priceHistAdapter := pricehistoryinterface.NewAdapter(
		*apppricehistory.NewCreatePriceHistoryUseCase(priceHistRepo),
		*apppricehistory.NewGetPriceHistoryUseCase(priceHistRepo),
	)

	productBundleAdapter := productbundleinterface.NewAdapter(
		*appproductbundle.NewCreateBundleUseCase(productBundleRepo),
		*appproductbundle.NewListBundlesUseCase(productBundleRepo),
	)

	refPriceAdapter := referencepriceinterface.NewAdapter(
		appreferenceprice.NewCreateReferencePriceUseCase(refPriceRepo),
		appreferenceprice.NewGetReferencePriceUseCase(refPriceRepo),
		appreferenceprice.NewGetByProductReferencePriceUseCase(refPriceRepo),
		appreferenceprice.NewListReferencePricesUseCase(refPriceRepo),
		appreferenceprice.NewDeleteReferencePriceUseCase(refPriceRepo),
		appreferenceprice.NewValidateReferencePriceUseCase(refPriceRepo, inventoryRepo),
	)

	salesCommAdapter := salescommissioninterface.NewAdapter(
		appsalescommission.NewCreateSalesCommissionUseCase(salesCommRepo),
		appsalescommission.NewUpdateMaxPriceUseCase(salesCommRepo),
		appsalescommission.NewUpdateMinQtyUseCase(salesCommRepo),
		appsalescommission.NewGetSalesCommissionUseCase(salesCommRepo),
		appsalescommission.NewGetByInventorySalesCommissionUseCase(salesCommRepo),
		appsalescommission.NewListSalesCommissionsUseCase(salesCommRepo),
		appsalescommission.NewDeleteSalesCommissionUseCase(salesCommRepo),
		appsalescommission.NewCalculateCommissionUseCase(salesCommRepo, inventoryRepo),
	)

	catCommRuleAdapter := salescommissioninterface.NewCategoryCommissionRuleAdapter(
		appsalescommission.NewCreateCategoryCommissionRuleUseCase(catCommRuleRepo),
		appsalescommission.NewGetCategoryCommissionRuleUseCase(catCommRuleRepo),
		appsalescommission.NewListCategoryCommissionRulesUseCase(catCommRuleRepo),
		appsalescommission.NewUpdateCategoryCommissionRuleUseCase(catCommRuleRepo),
		appsalescommission.NewDeleteCategoryCommissionRuleUseCase(catCommRuleRepo),
	)

	storeCatAdapter := storeallowedcategoryinterface.NewAdapter(
		appstoreallowedcategory.NewCreateCategoryUseCase(storeCatRepo),
		appstoreallowedcategory.NewApproveCategoryUseCase(storeCatRepo),
		appstoreallowedcategory.NewRejectCategoryUseCase(storeCatRepo),
	)

	createLinkUC := appstorewarehouselink.NewCreateLinkUseCase(warehouseLinkRepo)
	changeRelationUC := appstorewarehouselink.NewChangeRelationUseCase(warehouseLinkRepo)

	warehouseLinkAdapter := storewarehouselinkinterface.NewAdapter(
		&createLinkAdapter{inner: createLinkUC},
		&changeRelationAdapter{uc: changeRelationUC, repo: warehouseLinkRepo},
	)

	createWHUC := appwarehouse.NewCreateWarehouseUseCase(warehouseRepo)
	updateVisUC := appwarehouse.NewUpdateVisibilityUseCase(warehouseRepo)
	updateContUC := appwarehouse.NewUpdateContactUseCase(warehouseRepo)

	warehouseAdapter := warehouseinterface.NewAdapter(
		&createWarehouseAdapter{inner: createWHUC},
		&updateVisibilityAdapter{uc: updateVisUC, repo: warehouseRepo},
		&updateContactAdapter{uc: updateContUC, repo: warehouseRepo},
	)

	mux := router.New(router.Config{
		Brand:            brandAdapter,
		Category:         categoryAdapter,
		Store:            storeAdapter,
		Inventory:        inventoryAdapter,
		Product:          productAdapter,
		ProductImage:     productImageAdapter,
		ProductType:      productTypeAdapter,
		ProductAttribute: productAttrAdapter,
		PriceHistory:     priceHistAdapter,
		ProductBundle:    productBundleAdapter,
		Promotion:        promotionAdapter,
		ReferencePrice:   refPriceAdapter,
		SalesCommission:  salesCommAdapter,
		CommissionRule:    catCommRuleAdapter,
		StoreCategory:    storeCatAdapter,
		WarehouseLink:    warehouseLinkAdapter,
		Warehouse:        warehouseAdapter,
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type createWarehouseAdapter struct {
	inner *appwarehouse.CreateWarehouseUseCase
}

func (a *createWarehouseAdapter) Execute(input warehouseinterface.CreateWarehouseInput) (*warehousedomain.Warehouse, error) {
	return a.inner.Execute(input.CreatedByUserID, input.WarehouseName)
}

type updateVisibilityAdapter struct {
	uc   *appwarehouse.UpdateVisibilityUseCase
	repo *memory.WarehouseRepository
}

func (a *updateVisibilityAdapter) Execute(input warehouseinterface.UpdateVisibilityInput) (*warehousedomain.Warehouse, error) {
	if err := a.uc.Execute(input.WarehouseID, input.IsPublic); err != nil {
		return nil, err
	}
	return a.repo.FindByID(input.WarehouseID)
}

type updateContactAdapter struct {
	uc   *appwarehouse.UpdateContactUseCase
	repo *memory.WarehouseRepository
}

func (a *updateContactAdapter) Execute(input warehouseinterface.UpdateContactInput) (*warehousedomain.Warehouse, error) {
	if err := a.uc.Execute(input.WarehouseID, input.Phone, input.ContactPhone, input.CollectionMethod); err != nil {
		return nil, err
	}
	return a.repo.FindByID(input.WarehouseID)
}

type createLinkAdapter struct {
	inner *appstorewarehouselink.CreateLinkUseCase
}

func (a *createLinkAdapter) Execute(input storewarehouselinkinterface.CreateLinkInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
	return a.inner.Execute(input.StoreID, input.WarehouseID)
}

type changeRelationAdapter struct {
	uc   *appstorewarehouselink.ChangeRelationUseCase
	repo *memory.WarehouseLinkRepository
}

func (a *changeRelationAdapter) Execute(input storewarehouselinkinterface.ChangeRelationInput) (*storewarehouselinkdomain.StoreWarehouseLink, error) {
	if err := a.uc.Execute(input.LinkID, storewarehouselinkdomain.RelationType(input.RelationType)); err != nil {
		return nil, err
	}
	return a.repo.FindByID(input.LinkID)
}
