package memory

import (
	"time"

	branddomain "stock-service/internal/domain/brand"
	categorydomain "stock-service/internal/domain/category"
	inventorydomain "stock-service/internal/domain/inventory"
	productdomain "stock-service/internal/domain/product"
	promotiondomain "stock-service/internal/domain/promotion"
	referencepricedomain "stock-service/internal/domain/reference_price"
	salescommissiondomain "stock-service/internal/domain/sales_commission"
	storedomain "stock-service/internal/domain/store"
	storeallowedcategorydomain "stock-service/internal/domain/store_allowed_category"
	storewarehouselinkdomain "stock-service/internal/domain/store_warehouse_link"
	warehousedomain "stock-service/internal/domain/warehouse"
)

func SeedAll(
	brandRepo *BrandRepository,
	categoryRepo *CategoryRepository,
	storeRepo *StoreRepository,
	warehouseRepo *WarehouseRepository,
	productRepo *ProductRepository,
	inventoryRepo *InventoryRepository,
	promotionRepo *PromotionRepository,
	refPriceRepo *ReferencePriceRepository,
	salesCommRepo *SalesCommissionRepository,
	catCommRuleRepo *CategoryCommissionRuleRepository,
	storeCatRepo *StoreCategoryRepository,
	warehouseLinkRepo *WarehouseLinkRepository,
) {
	seedBrands(brandRepo)
	seedCategories(categoryRepo)
	seedStores(storeRepo)
	seedWarehouses(warehouseRepo)
	seedProducts(productRepo)
	seedInventories(inventoryRepo, storeRepo, warehouseRepo, productRepo)
	seedPromotions(promotionRepo)
	seedReferencePrices(refPriceRepo, productRepo)
	seedCategoryCommissionRules(catCommRuleRepo, categoryRepo)
	seedSalesCommissions(salesCommRepo, inventoryRepo, catCommRuleRepo)
	seedStoreCategories(storeCatRepo, storeRepo, categoryRepo)
	seedWarehouseLinks(warehouseLinkRepo, storeRepo, warehouseRepo)
}

func seedBrands(r *BrandRepository) {
	brands := []struct {
		id   int64
		name string
		slug string
	}{
		{1, "سامسونگ", "samsung"},
		{2, "اپل", "apple"},
		{3, "شیائومی", "xiaomi"},
		{4, "هوآوی", "huawei"},
		{5, "نوکیا", "nokia"},
	}
	for _, b := range brands {
		brand, err := branddomain.NewBrand(b.name, b.slug)
		if err != nil {
			panic(err)
		}
		brand.ID = b.id
		r.items[brand.ID] = brand
	}
}

func seedCategories(r *CategoryRepository) {
	cats := []struct {
		id       int64
		name     string
		slug     string
		parentID *int64
	}{
		{1, "موبایل", "mobile", nil},
		{2, "لپ تاپ", "laptop", nil},
		{3, "گوشی هوشمند", "smartphone", int64Ptr(1)},
		{4, "تبلت", "tablet", int64Ptr(1)},
		{5, "لوازم جانبی", "accessories", nil},
	}
	for _, c := range cats {
		cat, err := categorydomain.NewCategory(c.name, c.slug, c.parentID)
		if err != nil {
			panic(err)
		}
		cat.ID = c.id
		r.items[cat.ID] = cat
	}
}

func seedStores(r *StoreRepository) {
	stores := []struct {
		id        int64
		userID    int64
		storeName string
	}{
		{1, 101, "فروشگاه دیجیتال"},
		{2, 102, "فروشگاه موبایل سنتر"},
		{3, 103, "فروشگاه لب تاپ لند"},
	}
	for _, s := range stores {
		store, err := storedomain.NewStore(s.userID, s.storeName)
		if err != nil {
			panic(err)
		}
		store.ID = s.id
		r.stores[store.ID] = store
	}
}

func seedWarehouses(r *WarehouseRepository) {
	whs := []struct {
		id               int64
		createdByUserID  int64
		name             string
		collectionMethod string
		isPublic         bool
	}{
		{1, 1, "انبار مرکزی تهران", "pickup", true},
		{2, 1, "انبار غرب", "delivery", true},
		{3, 1, "انبار شمال", "both", false},
		{4, 2, "انبار شرق", "pickup", true},
	}
	for _, w := range whs {
		wh, err := warehousedomain.NewWarehouse(w.createdByUserID, w.name)
		if err != nil {
			panic(err)
		}
		wh.ID = w.id
		wh.CollectionMethod = w.collectionMethod
		if w.isPublic {
			wh.MakePublic()
		}
		r.items[wh.ID] = wh
	}
}

func seedProducts(r *ProductRepository) {
	type seedProduct struct {
		id         int32
		titleFa    string
		slug       string
		brandID    int64
		categoryID int64
		isEnabled  bool
		status     productdomain.ProductStatus
	}
	products := []seedProduct{
		{1, "محصول یک", "product-1", 1, 3, true, productdomain.ProductStatusActive},
		{2, "محصول دو", "product-2", 2, 3, true, productdomain.ProductStatusActive},
		{3, "محصول سه", "product-3", 3, 3, false, productdomain.ProductStatusPending},
		{10, "محصول ده", "product-10", 1, 2, true, productdomain.ProductStatusActive},
		{30, "محصول سی", "product-30", 2, 2, true, productdomain.ProductStatusActive},
		{42, "محصول چهل و دو", "product-42", 3, 4, false, productdomain.ProductStatusActive},
		{100, "محصول صد", "product-100", 1, 5, true, productdomain.ProductStatusActive},
		{200, "محصول دویست", "product-200", 4, 3, true, productdomain.ProductStatusActive},
		{300, "محصول سیصد", "product-300", 5, 5, false, productdomain.ProductStatusPending},
	}
	for _, p := range products {
		product, err := productdomain.NewProduct(p.titleFa, p.brandID, p.categoryID, productdomain.WithSlug(p.slug))
		if err != nil {
			panic(err)
		}
		product.ID = p.id
		product.Status = p.status
		if p.isEnabled {
			product.MarkEnabled()
		}
		r.items[product.ID] = product
	}
}

func seedInventories(invRepo *InventoryRepository, storeRepo *StoreRepository, whRepo *WarehouseRepository, prodRepo *ProductRepository) {
	invRepo.mu.Lock()
	defer invRepo.mu.Unlock()

	types := []struct {
		id          int64
		storeID     int64
		warehouseID int64
		productID   int32
		basePrice   float64
		instantQty  int
		// optionally apply a promotion
		promotionID *int64
		finalPrice  *float64
	}{
		{1, 1, 1, 1, 1500000, 50, nil, nil},
		{2, 1, 1, 2, 2500000, 30, int64Ptr(1), float64Ptr(2200000)},
		{3, 1, 2, 3, 1200000, 100, nil, nil},
		{4, 2, 1, 10, 35000000, 10, nil, nil},
		{5, 2, 3, 30, 45000000, 5, nil, nil},
		{6, 3, 4, 42, 800000, 200, nil, nil},
		{7, 1, 4, 100, 50000, 1000, int64Ptr(2), float64Ptr(40000)},
		{8, 2, 2, 200, 1800000, 75, nil, nil},
		{9, 3, 1, 300, 300000, 500, nil, nil},
	}
	for _, t := range types {
		inv, err := inventorydomain.NewInventory(t.storeID, t.warehouseID, t.productID, t.basePrice)
		if err != nil {
			panic(err)
		}
		inv.ID = t.id
		inv.InstantQty = t.instantQty
		if t.promotionID != nil {
			inv.PromotionID = t.promotionID
			inv.FinalPrice = t.finalPrice
		}
		invRepo.items[inv.ID] = inv
	}
}

func seedPromotions(r *PromotionRepository) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now.Add(7 * 24 * time.Hour)

	promos := []promotiondomain.CreatePromotionInput{
		{
			Title:         "تخفیف ۱۰ درصدی",
			DiscountType:  promotiondomain.DiscountTypePercentage,
			DiscountValue: 10,
			StartAt:       &start,
			EndAt:         &end,
		},
		{
			Title:         "۵۰ هزار تومان تخفیف",
			DiscountType:  promotiondomain.DiscountTypeFixedAmount,
			DiscountValue: 50000,
			StartAt:       &start,
			EndAt:         &end,
		},
		{
			Title:         "تخفیف ۲۰ درصدی با کد WELCOME",
			DiscountType:  promotiondomain.DiscountTypePercentage,
			DiscountValue: 20,
			CouponCode:    strPtr("WELCOME"),
			UsageLimit:    intPtr(100),
			StartAt:       &start,
			EndAt:         &end,
		},
		{
			Title:         "تخفیف ویژه وسایل جانبی",
			DiscountType:  promotiondomain.DiscountTypePercentage,
			DiscountValue: 15,
			EligibleCategoryIDs: []int64{5},
			StartAt:       &start,
			EndAt:         &end,
		},
	}
	for i, input := range promos {
		p, err := promotiondomain.NewPromotion(input)
		if err != nil {
			panic(err)
		}
		p.ID = int64(i + 1)
		p.Activate()
		p.UsedCount = 5
		if p.Budget != nil {
			p.BudgetSpent = *p.Budget * 0.3
		}
		r.items[p.ID] = p
	}
}

func seedReferencePrices(r *ReferencePriceRepository, prodRepo *ProductRepository) {
	refs := []struct {
		productID int32
		price     float64
		source    string
	}{
		{1, 1600000, "digikala"},
		{2, 2700000, "torob"},
		{3, 1300000, "digikala"},
		{10, 37000000, "emalls"},
		{100, 55000, "digikala"},
	}
	for i, ref := range refs {
		rp, err := referencepricedomain.NewReferencePrice(ref.productID, ref.price, ref.source)
		if err != nil {
			panic(err)
		}
		rp.ID = int64(i + 1)
		r.items[rp.ID] = rp
	}
}

func seedCategoryCommissionRules(r *CategoryCommissionRuleRepository, catRepo *CategoryRepository) {
	rules := []struct {
		categoryID  int32
		ratePercent float64
		minPrice    float64
		isActive    bool
	}{
		{1, 5.0, 100000, true},
		{2, 3.0, 500000, true},
		{3, 4.5, 200000, true},
		{5, 2.0, 50000, true},
	}
	for i, rule := range rules {
		cr, err := salescommissiondomain.NewCategoryCommissionRule(rule.categoryID, rule.ratePercent, rule.minPrice)
		if err != nil {
			panic(err)
		}
		cr.ID = int64(i + 1)
		cr.IsActive = rule.isActive
		r.items[cr.ID] = cr
	}
}

func seedSalesCommissions(salesCommRepo *SalesCommissionRepository, invRepo *InventoryRepository, catCommRuleRepo *CategoryCommissionRuleRepository) {
	comms := []struct {
		inventoryID int64
		ruleID      int64
		ratePercent float64
		minPrice    float64
	}{
		{1, 1, 5.0, 100000},
		{2, 1, 5.0, 100000},
		{3, 2, 3.5, 200000},
		{4, 2, 3.5, 200000},
		{5, 3, 4.5, 500000},
	}
	for i, c := range comms {
		sc, err := salescommissiondomain.NewSalesCommission(c.inventoryID, c.ruleID, salescommissiondomain.SaleModelRetail, c.ratePercent, c.minPrice)
		if err != nil {
			panic(err)
		}
		sc.ID = int64(i + 1)
		salesCommRepo.items[sc.ID] = sc
	}
}

func seedStoreCategories(r *StoreCategoryRepository, storeRepo *StoreRepository, catRepo *CategoryRepository) {
	scs := []struct {
		storeID    int64
		categoryID int64
		status     storeallowedcategorydomain.Status
	}{
		{1, 3, storeallowedcategorydomain.StatusApproved},
		{1, 5, storeallowedcategorydomain.StatusApproved},
		{2, 3, storeallowedcategorydomain.StatusApproved},
		{2, 2, storeallowedcategorydomain.StatusPending},
		{3, 2, storeallowedcategorydomain.StatusApproved},
	}
	for i, sc := range scs {
		sac := storeallowedcategorydomain.NewStoreAllowedCategory(sc.storeID, sc.categoryID)
		sac.ID = int64(i + 1)
		if sc.status == storeallowedcategorydomain.StatusApproved {
			sac.Approve()
		}
		r.items[sac.ID] = sac
	}
}

func seedWarehouseLinks(r *WarehouseLinkRepository, storeRepo *StoreRepository, whRepo *WarehouseRepository) {
	links := []struct {
		storeID     int64
		warehouseID int64
	}{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 3},
		{3, 4},
		{3, 1},
	}
	for i, l := range links {
		swl := storewarehouselinkdomain.NewStoreWarehouseLink(l.storeID, l.warehouseID)
		swl.ID = int64(i + 1)
		r.items[swl.ID] = swl
	}
}

func int64Ptr(v int64) *int64  { return &v }
func strPtr(v string) *string  { return &v }
func float64Ptr(v float64) *float64 { return &v }
func intPtr(v int) *int        { return &v }
