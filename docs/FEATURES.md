# Stock Service — Feature Inventory

> All 9 domains. Every entity, use case, endpoint, and test.
> ✅ = fully implemented, 🔶 = partial/stub, ❌ = not implemented

---

## 1. Store

**Domain entity** (`internal/domain/store/store.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `UserID` | `int64` | ✅ |
| `StoreName` | `string` | ✅ |
| `Status` | `StoreStatus` ("active") | ✅ |
| `AddressID` | `*int64` | ✅ (stored, no behavior) |
| `ContactPhone` | `*string` | ✅ |
| `MediaAssets` | `map[string]any` | ✅ (stored, no behavior) |
| `IsCommissionApplicable` | `bool` | ✅ |
| `IsBulkSaleEnabled` | `bool` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewStore(userID, storeName)` | ✅ validates name, defaults active/commission/no-bulk |
| `EnableBulkSale()` / `DisableBulkSale()` | ✅ |
| `EnableCommission()` / `DisableCommission()` | ✅ |
| `UpdateName(name)` | ✅ validates name |
| `UpdateContactInfo(phone)` | ✅ |

**Repository interface** (`internal/domain/store/repository.go`)

| Method | Status |
|---|---|
| `Save(store)` | ✅ |
| `FindByID(id)` | ✅ |
| `FindAll()` | ✅ |
| `Delete(id)` | ✅ |

**Validators & errors**

| Validator | Error | Status |
|---|---|---|
| `ValidateStoreName` | `ErrStoreNameRequired`, `ErrStoreNameTooLong` | ✅ |

**Use cases** (`internal/application/store/`)

| Use Case | Signature | Status |
|---|---|---|
| CreateStore | `Execute(CreateStoreInput) (*Store, error)` | ✅ |
| GetStore | `Execute(GetStoreInput) (*Store, error)` | ✅ |
| ListStores | `Execute(ListStoresInput) (*ListStoresOutput, error)` | ✅ filters by user, status; paginated |
| DeleteStore | `Execute(DeleteStoreInput) error` | ✅ |
| UpdateStoreName | `Execute(UpdateStoreNameInput) (*Store, error)` | ✅ validates name |
| UpdateStoreProfile | `Execute(UpdateStoreProfileInput) (*Store, error)` | ✅ updates address_id + media_assets |
| ToggleBulkSale | `Execute(ToggleBulkSaleInput) (*Store, error)` | ✅ |
| ToggleCommission | `Execute(ToggleCommissionInput) (*Store, error)` | ✅ |
| UpdateContact | `Execute(UpdateContactInput) (*Store, error)` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/stores` | POST | ✅ |
| `/api/v1/stores` | GET | ✅ list (user_id, status, page, limit) |
| `/api/v1/stores/{id}` | GET | ✅ |
| `/api/v1/stores/{id}` | PUT | ✅ update name |
| `/api/v1/stores/{id}` | DELETE | ✅ |
| `/api/v1/stores/{id}/profile` | PUT | ✅ update address + media assets |
| `/api/v1/stores/{id}/bulk-sale` | POST | ✅ |
| `/api/v1/stores/{id}/commission` | POST | ✅ |
| `/api/v1/stores/{id}/contact` | PUT | ✅ |

**Missing Store features**

| Feature | Status |
|---|---|
| Store-level permissions/vendor gating | ❌ (needs auth layer) |

---

## 2. Inventory

**Domain entity** (`internal/domain/inventory/inventory.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `StoreID` | `int64` | ✅ |
| `WarehouseID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `SaleModel` | `SaleModel` ("retail") | ✅ (only retail) |
| `BasePrice` | `float64` | ✅ |
| `PromotionID` | `*int64` | ✅ |
| `FinalPrice` | `*float64` | ✅ |
| `StartAt` / `EndAt` | `*time.Time` | ✅ |
| `PromotionStatus` | `PromotionStatus` ("pending") | ✅ |
| `Attributes` | `map[string]any` | ✅ (stored, no behavior) |
| `InstantQty` | `int` | ✅ |
| `ScheduledQty` | `map[string]int` | ✅ |
| `MinOrderQty` | `int` | ✅ |
| `MaxOrderQty` | `*int` | ✅ |
| `Condition` | `Condition` ("new") | ✅ |
| `VendorSaleStatus` | `VendorSaleStatus` ("active") | ✅ |
| `SystemSaleStatus` | `SystemSaleStatus` ("active") | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewInventory(storeID, warehouseID, productID, basePrice)` | ✅ validates base price > 0 |
| `ApplyPromotion(promotionID, finalPrice, startAt, endAt)` | ✅ |
| `RemovePromotion()` | ✅ |
| `UpdateInventory(instantQty, scheduledQty, minOrderQty, maxOrderQty)` | ✅ validates quantities |
| `SuspendVendorSale()` | ✅ |
| `CloseVendorSale()` | ✅ |
| `SuspendSystemSale()` | ✅ |
| `CloseSystemSale()` | ✅ |
| `ReserveQuantity(qty)` | ✅ |
| `ReleaseQuantity(qty)` | ✅ |
| `HasLowStock(threshold)` | ✅ |
| `ValidateScheduledQty(deliveryDate, qty)` | ✅ |

**Validators & errors**

| Validator | Error | Status |
|---|---|---|
| `ValidateBasePrice` | `ErrInvalidBasePrice` | ✅ |
| `ValidateFinalPrice` | `ErrInvalidFinalPrice` | ✅ |
| `ValidatePromotionDates` | `ErrInvalidPromotionDates` | ✅ |
| `ValidateInstantQty` | `ErrInvalidInstantQty` | ✅ |
| `ValidateMinOrderQty` | `ErrInvalidMinOrderQty` | ✅ |
| `ValidateMaxOrderQty` | `ErrInvalidMaxOrderQty` | ✅ |
| `ValidateScheduledDeliveryDate` | `ErrInvalidScheduledDate` | ✅ |
| — | `ErrPromotionAlreadyApplied` | ✅ |
| — | `ErrNoActivePromotion` | ✅ |
| — | `ErrProductNotFound` | ✅ (product existence check) |
| — | `ErrInventoryNotFound` | ✅ |
| — | `ErrVendorSaleStatusTransition` | ✅ |
| — | `ErrSystemSaleStatusTransition` | ✅ |
| — | `ErrInsufficientStock` | ✅ |

**Repository interface** (`internal/domain/inventory/repository.go`)

| Method | Status |
|---|---|
| `Save(inv)` | ✅ |
| `FindByID(id)` | ✅ |
| `FindAll()` | ✅ |
| `Delete(id)` | ✅ |

**Use cases** (`internal/application/inventory/`)

| Use Case | Signature | Status |
|---|---|---|
| CreateInventory | `Execute(CreateInventoryInput) (*Inventory, error)` | ✅ validates product exists |
| GetInventory | `Execute(GetInventoryInput) (*Inventory, error)` | ✅ |
| ListInventory | `Execute(ListInventoryInput) (*ListInventoryOutput, error)` | ✅ filters by store, product, vendor/system status; paginated |
| DeleteInventory | `Execute(DeleteInventoryInput) error` | ✅ |
| SearchInventory | `Execute(SearchInventoryInput) (*SearchInventoryOutput, error)` | ✅ by product name, paginated |
| ApplyPromotion | `Execute(ApplyPromotionInput) (*Inventory, error)` | ✅ |
| RemovePromotion | `Execute(RemovePromotionInput) (*Inventory, error)` | ✅ |
| UpdateInventory | `Execute(UpdateInventoryInput) (*Inventory, error)` | ✅ |
| SuspendVendorSale | `Execute(SuspendVendorSaleInput) (*Inventory, error)` | ✅ |
| CloseVendorSale | `Execute(CloseVendorSaleInput) (*Inventory, error)` | ✅ |
| SuspendSystemSale | `Execute(SuspendSystemSaleInput) (*Inventory, error)` | ✅ |
| CloseSystemSale | `Execute(CloseSystemSaleInput) (*Inventory, error)` | ✅ |
| ReserveQuantity | `Execute(ReserveQuantityInput) (*Inventory, error)` | ✅ |
| ReleaseQuantity | `Execute(ReleaseQuantityInput) (*Inventory, error)` | ✅ |
| CheckLowStock | `Execute(CheckLowStockInput) (*CheckLowStockOutput, error)` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/inventory` | POST | ✅ |
| `/api/v1/inventory` | GET | ✅ list (store_id, product_id, vendor_sale_status, system_sale_status, page, limit) |
| `/api/v1/inventory/search` | GET | ✅ search by product name |
| `/api/v1/inventory/{id}` | GET | ✅ |
| `/api/v1/inventory/{id}` | DELETE | ✅ |
| `/api/v1/inventory/{id}/promotion` | POST | ✅ |
| `/api/v1/inventory/{id}/promotion` | DELETE | ✅ |
| `/api/v1/inventory/{id}/inventory` | PUT | ✅ |
| `/api/v1/inventory/{id}/vendor/suspend` | POST | ✅ |
| `/api/v1/inventory/{id}/vendor/close` | POST | ✅ |
| `/api/v1/inventory/{id}/system/suspend` | POST | ✅ |
| `/api/v1/inventory/{id}/system/close` | POST | ✅ |
| `/api/v1/inventory/{id}/reserve` | POST | ✅ |
| `/api/v1/inventory/{id}/release` | POST | ✅ |
| `/api/v1/inventory/{id}/low-stock` | GET | ✅ |

**Missing Inventory features**

| Feature | Status |
|---|---|
| — all planned features implemented | ✅ |

---

## 3. Product

**Domain entity** (`internal/domain/product/product.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int32` | ✅ |
| `TitleFa` | `string` | ✅ |
| `TitleEn` | `*string` | ✅ |
| `Description` | `*string` | ✅ |
| `Slug` | `string` | ✅ |
| `MetaTitle` | `*string` | ✅ |
| `MetaDescription` | `*string` | ✅ |
| `IsEnabled` | `bool` | ✅ |
| `EnabledAt` | `*time.Time` | ✅ |
| `DisabledAt` | `*time.Time` | ✅ |
| `BrandID` | `int64` | ✅ |
| `CategoryID` | `int64` | ✅ |
| `OwnerType` | `OwnerType` (system/user) | ✅ |
| `OwnerID` | `*int64` | ✅ |
| `IsOriginal` | `bool` | ✅ (defaults true) |
| `Status` | `ProductStatus` (pending/active/rejected/deleted) | ✅ |
| `CreatedAt` | `time.Time` | ✅ |
| `UpdatedAt` | `time.Time` | ✅ |
| `IndexImageFileID` | `*int64` | ✅ |
| `DeletedAt` | `*time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewProduct(titleFa, brandID, categoryID, opts...)` | ✅ validates required, applies 6 functional options |
| `MarkActive()` | ✅ |
| `MarkRejected()` | ✅ |
| `MarkEnabled()` / `MarkDisabled()` | ✅ |
| `SoftDelete()` | ✅ sets status + DeletedAt |
| `UpdateSEO(metaTitle, metaDesc)` | ✅ |
| `GenerateSlug(slug)` | ✅ |
| `Touch()` | ✅ updates UpdatedAt |

**Functional options**

| Option | Status |
|---|---|
| `WithTitleEn(v)` | ✅ |
| `WithDescription(v)` | ✅ |
| `WithSlug(v)` | ✅ |
| `WithMetaTitle(v)` | ✅ |
| `WithMetaDescription(v)` | ✅ |
| `WithOwnerType(v)` | ✅ |
| `WithOwnerID(v)` | ✅ |
| `WithIsOriginal(v)` | ✅ |
| `WithIndexImageFileID(v)` | ✅ |

**Validators & errors**

| Validation | Status |
|---|---|
| title_fa required → `ErrTitleFaRequired` | ✅ |
| brand_id > 0 → `ErrInvalidBrandID` | ✅ |
| category_id > 0 → `ErrInvalidCategoryID` | ✅ |
| ID > 0 → `ErrInvalidProductID` | ✅ |
| not found → `ErrProductNotFound` | ✅ |

**Use cases** (`internal/application/product/`)

| Use Case | Signature | Status |
|---|---|---|
| CreateProduct | `Execute(CreateProductInput) (*Product, error)` | ✅ |
| GetProduct | `Execute(GetProductInput) (*Product, error)` | ✅ |
| UpdateProduct | `Execute(UpdateProductInput) (*Product, error)` | ✅ |
| ActivateProduct | `Execute(ActivateProductInput) (*Product, error)` | ✅ |
| RejectProduct | `Execute(RejectProductInput) (*Product, error)` | ✅ |
| SoftDeleteProduct | `Execute(SoftDeleteProductInput) (*Product, error)` | ✅ |
| EnableProduct | `Execute(EnableProductInput) (*Product, error)` | ✅ |
| DisableProduct | `Execute(DisableProductInput) (*Product, error)` | ✅ |
| UpdateSEO | `Execute(UpdateSEOInput) (*Product, error)` | ✅ |
| ListProducts | `Execute(ListProductsInput) (*ListProductsOutput, error)` | ✅ by owner, status, category, brand, search; paginated |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/products` | POST | ✅ |
| `/api/v1/products` | GET | ✅ list (owner_type, owner_id, status, category_id, brand_id, search, page, limit) |
| `/api/v1/products/my` | GET | ✅ seller self-service "my products" |
| `/api/v1/products/{id}` | GET | ✅ |
| `/api/v1/products/{id}` | PUT | ✅ |
| `/api/v1/products/{id}` | DELETE | ✅ soft delete |
| `/api/v1/products/{id}/activate` | POST | ✅ |
| `/api/v1/products/{id}/reject` | POST | ✅ |
| `/api/v1/products/{id}/enable` | POST | ✅ admin enable |
| `/api/v1/products/{id}/disable` | POST | ✅ admin disable |
| `/api/v1/products/{id}/seo` | PUT | ✅ update meta title / meta description |

**Brand entity** (`internal/domain/brand/`, `internal/application/brand/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `NameFa` | `string` | ✅ |
| `NameEn` | `*string` | ✅ |
| `Status` | `BrandStatus` (active/inactive) | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreateBrand, GetBrand, UpdateBrand, DeleteBrand, ListBrands — all ✅

**HTTP endpoints:** `POST/GET /api/v1/brands`, `GET/PUT/DELETE /api/v1/brands/{id}` — all ✅

**Category entity** (`internal/domain/category/`, `internal/application/category/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `NameFa` | `string` | ✅ |
| `NameEn` | `*string` | ✅ |
| `ParentID` | `*int64` | ✅ |
| `Status` | `CategoryStatus` (active/inactive) | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreateCategory, GetCategory, UpdateCategory, DeleteCategory, ListCategories — all ✅

**HTTP endpoints:** `POST/GET /api/v1/categories`, `GET/PUT/DELETE /api/v1/categories/{id}` — all ✅

**ProductImage** (`internal/domain/product_image/`, `internal/application/product_image/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `FileID` | `int64` | ✅ |
| `SortOrder` | `int` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreateImage, ListImages, DeleteImage — all ✅

**HTTP endpoints:** `POST/GET /api/v1/products/{productId}/images`, `DELETE /api/v1/products/images/{id}` — all ✅

**ProductType / Product Variants** (`internal/domain/product_type/`, `internal/application/product_type/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `Name` | `string` | ✅ (e.g. "Size", "Color") |
| `Value` | `string` | ✅ (e.g. "XL", "Red") |
| `SortOrder` | `int` | ✅ |

**Use cases:** CreateType, ListTypes — all ✅

**HTTP endpoints:** `POST/GET /api/v1/products/{productId}/types` — all ✅

**ProductAttribute / Custom Fields** (`internal/domain/product_attribute/`, `internal/application/product_attribute/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `Key` | `string` | ✅ |
| `Value` | `string` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreateAttribute, ListAttributes — all ✅

**HTTP endpoints:** `POST/GET /api/v1/products/{productId}/attributes` — all ✅

**PriceHistory** (`internal/domain/price_history/`, `internal/application/price_history/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `OldPrice` | `float64` | ✅ |
| `NewPrice` | `float64` | ✅ |
| `ChangedBy` | `string` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreatePriceHistory, GetPriceHistory — all ✅

**HTTP endpoints:** `POST/GET /api/v1/products/{productId}/price-history` — all ✅

**ProductBundle / Upsells / Cross-sells** (`internal/domain/product_bundle/`, `internal/application/product_bundle/`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `LinkedProductID` | `int64` | ✅ |
| `RelationType` | `RelationType` (bundle/upsell/cross-sell) | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Use cases:** CreateBundle, ListBundles — all ✅

**HTTP endpoints:** `POST/GET /api/v1/products/{productId}/bundles` — all ✅

**Missing Product features**

| Feature | Status |
|---|---|
| List products (by owner, by status, by category, by brand, search, paginated) | ✅ |
| Product images / gallery (multiple files, ordering) | ✅ |
| Product types / variants (size, color, etc.) | ✅ |
| Product attributes / custom fields | ✅ |
| Brand entity CRUD | ✅ |
| Category entity CRUD | ✅ |
| Seller self-service: "my products" list | ✅ |
| Admin enable/disable product | ✅ |
| Slug / URL-friendly name generation | ✅ |
| SEO fields (meta title, meta description) | ✅ |
| Price history | ✅ |
| Product bundles / upsells / cross-sells | ✅ |
| Review / rating system | ❌ |

---

## 4. Promotion

**Domain entity** (`internal/domain/promotion/promotion.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `Title` | `string` | ✅ |
| `DiscountType` | `DiscountType` (percentage/fixed_amount) | ✅ |
| `DiscountValue` | `float64` | ✅ |
| `MinPurchase` | `*float64` | ✅ |
| `CouponCode` | `*string` | ✅ |
| `UsageLimit` | `*int` | ✅ |
| `UsedCount` | `int` | ✅ |
| `MaxDiscountAmount` | `*float64` | ✅ percentage cap |
| `Budget` | `*float64` | ✅ |
| `BudgetSpent` | `float64` | ✅ |
| `EligibleStoreIDs` | `[]int64` | ✅ |
| `EligibleCategoryIDs` | `[]int64` | ✅ |
| `EligibleProductIDs` | `[]int32` | ✅ |
| `EligibleUserIDs` | `[]int64` | ✅ |
| `RequiresApproval` | `bool` | ✅ |
| `StartAt` | `*time.Time` | ✅ |
| `EndAt` | `*time.Time` | ✅ |
| `IsCountdown` | `bool` | ✅ |
| `ExpireSaleWithPromotion` | `bool` | ✅ |
| `Status` | `PromotionStatus` (inactive/active) | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewPromotion(input)` | ✅ validates title, discount type/value, dates, coupon code |
| `Update(input)` | ✅ partial update of any field with validation |
| `Activate()` | ✅ status → active |
| `Deactivate()` | ✅ status → inactive |
| `IsActive()` | ✅ |
| `IsExpired()` | ✅ checks EndAt |
| `IsScheduled()` | ✅ checks StartAt |
| `CanApply()` | ✅ active + not expired + within usage limit + within budget |
| `RecordUsage()` | ✅ increments UsedCount |
| `SpendBudget(amount)` | ✅ increments BudgetSpent |
| `CalculateDiscountPrice(basePrice)` | ✅ percentage (with optional cap) / fixed amount / min purchase |
| `IsEligibleForStore(id)` | ✅ empty list = no restriction |
| `IsEligibleForCategory(id)` | ✅ empty list = no restriction |
| `IsEligibleForProduct(id)` | ✅ empty list = no restriction |
| `IsEligibleForUser(id)` | ✅ empty list = no restriction |

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreatePromotion | `Execute(CreatePromotionInput) (*Promotion, error)` | ✅ full input with discount rules, eligibility, budget, schedule |
| GetPromotion | `Execute(GetPromotionInput) (*Promotion, error)` | ✅ |
| UpdatePromotion | `Execute(UpdatePromotionInput) (*Promotion, error)` | ✅ partial update |
| DeletePromotion | `Execute(DeletePromotionInput) error` | ✅ |
| ListPromotions | `Execute(ListPromotionsInput) (*ListPromotionsOutput, error)` | ✅ by status, discount type, search; paginated |
| ActivatePromotion | `Execute(id int64) error` | ✅ |
| DeactivatePromotion | `Execute(id int64) error` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/promotions` | POST | ✅ |
| `/api/v1/promotions` | GET | ✅ list (status, discount_type, search, page, limit) |
| `/api/v1/promotions/{id}` | GET | ✅ |
| `/api/v1/promotions/{id}` | PUT | ✅ update |
| `/api/v1/promotions/{id}` | DELETE | ✅ |
| `/api/v1/promotions/{id}/activate` | POST | ✅ |
| `/api/v1/promotions/{id}/deactivate` | POST | ✅ |

**Inventory auto-apply integration**

| Feature | Status |
|---|---|
| ApplyPromotion validates promotion exists, active, not expired, within budget/limit | ✅ |
| Auto-calculates FinalPrice from promotion discount rules when input is zero | ✅ |
| Checks store and product eligibility before applying | ✅ |
| Records usage count and budget spent on the promotion | ✅ |

**Missing Promotion features**

| Feature | Status |
|---|---|
| BOGO (buy one get one) discount type | ❌ |
| Tiered discount rules (e.g. 10% over $100, 15% over $200) | ❌ |
| Automatic promotion application on inventory creation | ❌ |
| Category and user eligibility checked during apply (stored but not validated) | 🔶 |

---

## 5. Reference Price

**Domain entity** (`internal/domain/reference_price/reference_price.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `ProductID` | `int32` | ✅ |
| `Price` | `float64` | ✅ |
| `Source` | `string` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewReferencePrice(productID, price, source)` | ✅ validates price > 0 |

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreateReferencePrice | `Execute(productID int32, price float64, source string) (*ReferencePrice, error)` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/reference-prices` | POST | ✅ |

**Missing ReferencePrice features**

| Feature | Status |
|---|---|
| Get reference price by ID | ❌ |
| Get reference price by product | ❌ |
| List/filter reference prices | ❌ |
| Delete reference price | ❌ |
| Compare with inventory base price for validation | ❌ |

---

## 6. Sales Commission

**Domain entity** (`internal/domain/sales_commission/sales_commission.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `InventoryID` | `int64` | ✅ |
| `CategoryCommissionRuleID` | `int64` | ✅ |
| `SaleModel` | `SaleModel` ("retail") | ✅ |
| `RatePercent` | `float64` | ✅ (0–100) |
| `MinQty` | `*int` | ✅ |
| `MinPrice` | `float64` | ✅ (> 0) |
| `MaxPrice` | `*float64` | ✅ (> min) |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewSalesCommission(inventoryID, ruleID, saleModel, rate, minPrice)` | ✅ validates rate + min price |
| `UpdateMaxPrice(maxPrice)` | ✅ validates > min |
| `UpdateMinQty(qty)` | ✅ validates >= 0 |

**⚠️ Known gap:** No `CalculateCommission(price, qty) float64` method exists. The rate is stored but never applied.

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreateSalesCommission | `Execute(inventoryID, ruleID, saleModel, rate, minPrice) (*SalesCommission, error)` | ✅ |
| UpdateMaxPrice | `Execute(commissionID, maxPrice) error` | ✅ |
| UpdateMinQty | `Execute(commissionID, minQty) error` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/sales-commissions` | POST | ✅ |
| `/api/v1/sales-commissions/{id}/max-price` | PUT | ✅ |
| `/api/v1/sales-commissions/{id}/min-qty` | PUT | ✅ |

**Missing SalesCommission features**

| Feature | Status |
|---|---|
| Calculate commission amount for a sale (rate × price × qty) | ❌ |
| Get commission by ID | ❌ |
| Get commission by inventory | ❌ |
| List/filter commissions | ❌ |
| Delete commission | ❌ |
| CategoryCommissionRule entity CRUD | ❌ |
| Decide commission base price (base vs final when promotion active) | ❌ |

---

## 7. Store Allowed Category

**Domain entity** (`internal/domain/store_allowed_category/store_allowed_category.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `StoreID` | `int64` | ✅ |
| `CategoryID` | `int64` | ✅ |
| `Status` | `Status` (pending/approved/rejected) | ✅ |
| `SupportNote` | `string` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewStoreAllowedCategory(storeID, categoryID)` | ✅ |
| `Approve()` | 🔶 (sets status, no validation) |
| `Reject()` | 🔶 (sets status, no validation) |

**⚠️ Known gap:** Validators and most errors are stubs with TODO comments.

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreateCategory | `Execute(storeID, categoryID int64) (*StoreAllowedCategory, error)` | ✅ |
| ApproveCategory | `Execute(id int64) error` | ✅ |
| RejectCategory | `Execute(id int64) error` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/store-categories` | POST | ✅ |
| `/api/v1/store-categories/{id}/approve` | POST | ✅ |
| `/api/v1/store-categories/{id}/reject` | POST | ✅ |

**Missing StoreAllowedCategory features**

| Feature | Status |
|---|---|
| Get by ID | ❌ |
| List by store | ❌ |
| Delete | ❌ |
| Add support note on reject | ❌ |
| Category → product category linkage | ❌ |

---

## 8. Store–Warehouse Link

**Domain entity** (`internal/domain/store_warehouse_link/store_warehouse_link.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `StoreID` | `int64` | ✅ |
| `WarehouseID` | `int64` | ✅ |
| `RelationType` | `RelationType` ("primary") | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewStoreWarehouseLink(storeID, warehouseID)` | ✅ |
| `ChangeRelationType(relationType)` | ✅ |

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreateLink | `Execute(storeID, warehouseID) (*StoreWarehouseLink, error)` | ✅ |
| ChangeRelation | `Execute(linkID, relationType) error` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/warehouse-links` | POST | ✅ |
| `/api/v1/warehouse-links/{id}/relation` | PUT | ✅ |

**Missing StoreWarehouseLink features**

| Feature | Status |
|---|---|
| Get link by ID | ❌ |
| List links by store / by warehouse | ❌ |
| Delete link | ❌ |
| Validator for RelationType values | ❌ |

---

## 9. Warehouse

**Domain entity** (`internal/domain/warehouse/warehouse.go`)

| Field | Type | Status |
|---|---|---|
| `ID` | `int64` | ✅ |
| `CreatedByUserID` | `int64` | ✅ |
| `WarehouseName` | `string` | ✅ |
| `AddressID` | `*int64` | ✅ (stored, unused) |
| `Phone` | `*string` | ✅ |
| `ContactPhone` | `*string` | ✅ |
| `IsPublic` | `bool` | ✅ |
| `CollectionMethod` | `string` | ✅ |
| `CreatedAt` | `time.Time` | ✅ |

**Domain methods**

| Method | Status |
|---|---|
| `NewWarehouse(createdByUserID, warehouseName)` | ✅ validates name |
| `MakePublic()` / `MakePrivate()` | ✅ |
| `UpdatePhone(phone)` | ✅ |
| `UpdateContactPhone(phone)` | ✅ |
| `UpdateCollectionMethod(method)` | ✅ |

**Use cases**

| Use Case | Signature | Status |
|---|---|---|
| CreateWarehouse | `Execute(createdByUserID, warehouseName) (*Warehouse, error)` | ✅ |
| UpdateVisibility | `Execute(id, isPublic) error` | ✅ |
| UpdateContact | `Execute(id, phone, contactPhone, collectionMethod) error` | ✅ |

**HTTP endpoints**

| Route | Method | Status |
|---|---|---|
| `/api/v1/warehouses` | POST | ✅ |
| `/api/v1/warehouses/{id}/visibility` | PUT | ✅ |
| `/api/v1/warehouses/{id}/contact` | PUT | ✅ |

**Missing Warehouse features**

| Feature | Status |
|---|---|
| Get warehouse by ID | ❌ |
| List warehouses (by user, by visibility) | ❌ |
| Delete warehouse | ❌ |
| Update warehouse name, address | ❌ |
| Validate collection method enum values | ❌ |

---

## Cross-Cutting Gaps

| Feature | Status |
|---|---|
| **Auth / permissions** — any caller can call any endpoint | ❌ |
| **List/search endpoints** — Store, Inventory, Product, Brand, Category have list with filtering + pagination | 🔶 (5/9 domains) |
| **Pagination** — Store, Inventory, Product use cases support pagination | 🔶 |
| **Order / checkout** — entirely absent | ❌ |
| **User entity** — referenced via `UserID`, `OwnerID`, `CreatedByUserID` but no User domain exists | ❌ |
| **Address entity** — referenced via `AddressID` but no Address domain exists | ❌ |
| **Media / file entity** — referenced via `IndexImageFileID`, `MediaAssets` but no File domain exists | ❌ |
| **Category entity** — referenced via `CategoryID` in Product and StoreAllowedCategory | ✅ |
| **Brand entity** — referenced via `BrandID` in Product | ✅ |
| **Database** — all repos are in-memory, no PostgreSQL implementation | ❌ |
| **Commission calculation** — rate stored but never computed | ❌ |
| **Promotion usage tracking** — UsedCount and BudgetSpent tracked on apply | ✅ |
| **Promotion discount rules** — percentage, fixed amount, coupon, budget, eligibility | ✅ |
| **Global search** — no product search by title_fa/title_en | ❌ |
| **Validation on missing reference entities** — no FK integrity check (e.g., product must have a valid brand) | ❌ |

---

## Test Coverage

| Package | Files | Tests | Status |
|---|---|---|---|---|---|
| `tests/entity/*` | 8 files | Creation, validation errors, state transitions | ✅ |
| `tests/application/*` | 47 files | Every use case (success + error) | ✅ |
| `tests/interface/*` | 18 files | Every adapter method (success + error mapping) | ✅ |
| `tests/interface/http/*` | 9 files | Every endpoint (success + invalid JSON, invalid ID, errors) | ✅ |
| **Total** | **73 files** | **35 test suites** | **✅ all pass** |
