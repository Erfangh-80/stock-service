package main

import (
	"log"
	"net/http"

	appinventory "stock-service/internal/application/inventory"
	appproduct "stock-service/internal/application/product"
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

	inventoryinterface "stock-service/internal/interface/inventory"
	productinterface "stock-service/internal/interface/product"
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
	storeRepo := memory.NewStoreRepository()
	inventoryRepo := memory.NewInventoryRepository()
	promotionRepo := memory.NewPromotionRepository()
	warehouseRepo := memory.NewWarehouseRepository()
	refPriceRepo := memory.NewReferencePriceRepository()
	salesCommRepo := memory.NewSalesCommissionRepository()
	storeCatRepo := memory.NewStoreCategoryRepository()
	warehouseLinkRepo := memory.NewWarehouseLinkRepository()
	productRepo := memory.NewProductRepository()
	memory.SeedProducts(productRepo)

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
		appinventory.NewApplyPromotionUseCase(inventoryRepo),
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
	)

	productAdapter := productinterface.NewAdapter(
		appproduct.NewCreateProductUseCase(productRepo),
		appproduct.NewGetProductUseCase(productRepo),
		appproduct.NewUpdateProductUseCase(productRepo),
		appproduct.NewActivateProductUseCase(productRepo),
		appproduct.NewRejectProductUseCase(productRepo),
		appproduct.NewSoftDeleteProductUseCase(productRepo),
	)

	refPriceAdapter := referencepriceinterface.NewAdapter(
		appreferenceprice.NewCreateReferencePriceUseCase(refPriceRepo),
	)

	salesCommAdapter := salescommissioninterface.NewAdapter(
		appsalescommission.NewCreateSalesCommissionUseCase(salesCommRepo),
		appsalescommission.NewUpdateMaxPriceUseCase(salesCommRepo),
		appsalescommission.NewUpdateMinQtyUseCase(salesCommRepo),
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
		Store:           storeAdapter,
		Inventory:       inventoryAdapter,
		Product:         productAdapter,
		Promotion:       promotionAdapter,
		ReferencePrice:  refPriceAdapter,
		SalesCommission: salesCommAdapter,
		StoreCategory:   storeCatAdapter,
		WarehouseLink:   warehouseLinkAdapter,
		Warehouse:       warehouseAdapter,
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
