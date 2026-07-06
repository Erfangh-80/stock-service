package router

import (
	"net/http"

	"stock-service/internal/interface/http/handler"
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
)

type Config struct {
	Brand           *brandinterface.Adapter
	Category        *categoryinterface.Adapter
	Inventory       *inventoryinterface.Adapter
	Product         *productinterface.Adapter
	ProductImage    *productimageinterface.Adapter
	ProductType     *producttypeinterface.Adapter
	ProductAttribute *productattributeinterface.Adapter
	PriceHistory    *pricehistoryinterface.Adapter
	ProductBundle   *productbundleinterface.Adapter
	Store           *storeinterface.Adapter
	Promotion       *promotioninterface.Adapter
	ReferencePrice  *referencepriceinterface.Adapter
	SalesCommission *salescommissioninterface.Adapter
	CommissionRule  *salescommissioninterface.CategoryCommissionRuleAdapter
	StoreCategory   *storeallowedcategoryinterface.Adapter
	WarehouseLink   *storewarehouselinkinterface.Adapter
	Warehouse       *warehouseinterface.Adapter
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()

	handler.NewBrandHandler(cfg.Brand).Register(mux)
	handler.NewCategoryHandler(cfg.Category).Register(mux)
	handler.NewInventoryHandler(cfg.Inventory).Register(mux)
	handler.NewProductHandler(cfg.Product).Register(mux)
	handler.NewProductImageHandler(cfg.ProductImage).Register(mux)
	handler.NewProductTypeHandler(cfg.ProductType).Register(mux)
	handler.NewProductAttributeHandler(cfg.ProductAttribute).Register(mux)
	handler.NewPriceHistoryHandler(cfg.PriceHistory).Register(mux)
	handler.NewProductBundleHandler(cfg.ProductBundle).Register(mux)
	handler.NewStoreHandler(cfg.Store).Register(mux)
	handler.NewPromotionHandler(cfg.Promotion).Register(mux)
	handler.NewReferencePriceHandler(cfg.ReferencePrice).Register(mux)
	handler.NewSalesCommissionHandler(cfg.SalesCommission).Register(mux)
	handler.NewCategoryCommissionRuleHandler(cfg.CommissionRule).Register(mux)
	handler.NewStoreAllowedCategoryHandler(cfg.StoreCategory).Register(mux)
	handler.NewStoreWarehouseLinkHandler(cfg.WarehouseLink).Register(mux)
	handler.NewWarehouseHandler(cfg.Warehouse).Register(mux)

	mux.Handle("/swagger/", http.StripPrefix("/swagger", handler.NewSwaggerHandler()))

	return mux
}
