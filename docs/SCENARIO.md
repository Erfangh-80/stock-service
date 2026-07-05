# SCENARIO — Complete Use Case Catalog & End‑to‑End Flow

> Visual diagrams available in [`FLOW.mmd`](FLOW.mmd)

---

## All Implemented Use Cases (25 total)

### inventory
| Use Case | File | What it does |
|---|---|---|
| `CreateInventory` | `create_product_sale.go` | Creates a new inventory listing |
| `ApplyPromotionToSale` | `apply_promotion.go` | Attaches a promotion to an inventory item |
| `RemovePromotionFromSale` | `remove_promotion.go` | Removes promotion from an inventory item |
| `UpdateProductSaleInventory` | `update_inventory.go` | Updates stock quantities |

### store
| Use Case | File | What it does |
|---|---|---|
| `CreateStore` | `create_store.go` | Registers a new store |
| `ToggleBulkSale` | `toggle_bulk_sale.go` | Enables or disables bulk sale |
| `ToggleCommission` | `toggle_commission.go` | Enables or disables commission |
| `UpdateStoreContact` | `update_contact.go` | Updates store contact phone |

### promotion
| Use Case | File | What it does |
|---|---|---|
| `CreatePromotion` | `create_promotion.go` | Creates a new promotion campaign |
| `ActivatePromotion` | `activate_promotion.go` | Activates a promotion |
| `DeactivatePromotion` | `deactivate_promotion.go` | Deactivates a promotion |

### reference_price
| Use Case | File | What it does |
|---|---|---|
| `CreateReferencePrice` | `create_reference_price.go` | Records a reference price for a product |

### sales_commission
| Use Case | File | What it does |
|---|---|---|
| `CreateSalesCommission` | `create_sales_commission.go` | Sets commission rate for a sale |
| `UpdateMaxPrice` | `update_max_price.go` | Updates max price for commission |
| `UpdateMinQty` | `update_min_qty.go` | Updates min quantity for commission |

### store_allowed_category
| Use Case | File | What it does |
|---|---|---|
| `CreateCategory` | `create_category.go` | Adds a category to a store |
| `ApproveCategory` | `approve_category.go` | Approves a pending category |
| `RejectCategory` | `reject_category.go` | Rejects a pending category |

### store_warehouse_link
| Use Case | File | What it does |
|---|---|---|
| `CreateLink` | `create_link.go` | Links a store to a warehouse |
| `ChangeRelation` | `change_relation.go` | Changes the link relation type |

### warehouse
| Use Case | File | What it does |
|---|---|---|
| `CreateWarehouse` | `create_warehouse.go` | Creates a new warehouse |
| `UpdateVisibility` | `update_visibility.go` | Makes warehouse public or private |
| `UpdateContact` | `update_contact.go` | Updates warehouse contact info |

---

## End‑to‑End Scenario: "Launch a Discounted Product"

A seller wants to list a product in their store with a discount promotion.

### Use cases involved (in order)

```
1. CreateStore           ──  register the seller's store
2. CreateWarehouse       ──  set up a warehouse
3. CreateLink            ──  link store → warehouse
4. CreateInventory       ──  list the product for sale
5. ApplyPromotionToSale  ──  attach the discount
```

---

### Step 1 — Register the store

**Use case:** `CreateStore`

```go
// Application receives the input
interactor.Execute(CreateStoreInput{
    UserID:    1,
    StoreName: "Erfan's Electronics",
})
```

**What the interactor does:**
```
repo (in memory)
  │
  ├─ domain.NewStore(1, "Erfan's Electronics")
  │     ├─ validator.ValidateStoreName("Erfan's Electronics")  ✅
  │     └─ returns Store{UserID:1, StoreName:"Erfan's Electronics",
  │                      Status:"active", IsCommissionApplicable:true, ...}
  │
  ├─ repo.Save(store)
  │
  └─ returns output
```

**Result:** A store is created with ID = 1, status `active`, commission enabled.

---

### Step 2 — Create the warehouse

**Use case:** `CreateWarehouse`

```go
interactor.Execute(CreateWarehouseInput{
    CreatedByUserID: 1,
    WarehouseName:   "Main Warehouse",
})
```

**What the interactor does:**
```
  ├─ domain.NewWarehouse(1, "Main Warehouse")
  │     ├─ validator.ValidateWarehouseName("Main Warehouse")  ✅
  │     └─ returns Warehouse{ID:1, CreatedByUserID:1, Name:"Main Warehouse", ...}
  │
  ├─ repo.Save(warehouse)
  └─ returns output
```

**Result:** Warehouse created with ID = 1.

---

### Step 3 — Link the store to the warehouse

**Use case:** `CreateLink`

```go
interactor.Execute(CreateLinkInput{
    StoreID:     1,
    WarehouseID: 1,
})
```

**What the interactor does:**
```
  ├─ domain.NewStoreWarehouseLink(1, 1)
  │     └─ returns StoreWarehouseLink{StoreID:1, WarehouseID:1,
  │                                    RelationType:"primary"}
  │
  ├─ repo.Save(link)
  └─ returns output
```

**Result:** Store 1 is linked to Warehouse 1 as `primary`.

---

### Step 4 — List the product for sale

**Use case:** `CreateInventory`

```go
interactor.Execute(CreateInventoryInput{
    StoreID:     1,
    WarehouseID: 1,
    ProductID:   42,
    BasePrice:   100.00,
})
```

**What the interactor does:**
```
  ├─ domain.NewInventory(1, 1, 42, 100.00)
  │     ├─ validator.ValidateBasePrice(100.00)  ✅
  │     └─ returns Inventory{
  │            ID:1, StoreID:1, WarehouseID:1, ProductID:42,
  │            BasePrice:100.00, SaleModel:"retail",
  │            MinOrderQty:1, Condition:"new",
  │            VendorSaleStatus:"active", SystemSaleStatus:"active",
  │            PromotionStatus:"pending", ...
  │          }
  │
  ├─ repo.Save(inv)
  └─ returns output
```

**Result:** Product 42 is listed at $100.00 in Store 1 / Warehouse 1.

---

### Step 5 — Apply the promotion to the product

**Use case:** `ApplyPromotionToSale`

```go
interactor.Execute(ApplyPromotionInput{
    SaleID:      1,
    PromotionID: 1,
    FinalPrice:  89.99,
    StartAt:     "2026-07-10T00:00:00Z",
    EndAt:       "2026-07-20T00:00:00Z",
})
```

**What the interactor does:**
```
  ├─ repo.FindByID(1)  →  finds Inventory 1
  │
  ├─ inv.ApplyPromotion(1, 89.99, start, end)
  │     ├─ Already has promotion?            →  no  ✅
  │     ├─ endAt after startAt?              →  yes ✅
  │     ├─ finalPrice > 0?                   →  yes ✅
  │     └─ updates:
  │          PromotionID:     1
  │          FinalPrice:      89.99
  │          StartAt:         2026-07-10
  │          EndAt:           2026-07-20
  │          PromotionStatus: "pending"
  │
  ├─ repo.Save(sale)
  └─ returns output
```

**Result:** Product 42 now has a $10.01 discount from July 10–20.

---

## Final state of all entities after the scenario

```
Store 1:            "Erfan's Electronics"  (active, commission enabled)
Warehouse 1:        "Main Warehouse"       (private)
Link:               Store 1 ↔ Warehouse 1  (primary)
Inventory 1:        Product 42 @ $100      (discounted to $89.99, promo pending)
```

## What each layer owned

| Step | Application role | Domain role |
|---|---|---|
| 1 | Called `NewStore`, called `Save` | Validated store name, set defaults |
| 2 | Called `NewWarehouse`, called `Save` | Validated warehouse name |
| 3 | Called `NewStoreWarehouseLink`, called `Save` | Set default relation type |
| 4 | Called `NewInventory`, called `Save` | Validated price, set all defaults |
| 5 | Loaded sale, called `ApplyPromotion`, saved | Validated 3 business rules, modified state |

The application layer never made a single business decision — it only
orchestrated loading, calling domain methods, and persisting.
