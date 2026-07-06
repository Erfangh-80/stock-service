# Stock Service

A clean architecture (layered) stock and inventory management service written in Go.

## Architecture

```
internal/
├── domain/          # Enterprise business rules (pure Go, zero deps)
│   ├── brand/
│   ├── category/
│   ├── inventory/
│   ├── price_history/
│   ├── product/
│   ├── product_attribute/
│   ├── product_bundle/
│   ├── product_image/
│   ├── product_type/
│   ├── promotion/
│   ├── reference_price/
│   ├── sales_commission/
│   ├── store/
│   ├── store_allowed_category/
│   ├── store_warehouse_link/
│   └── warehouse/
├── application/     # Application use cases (depends only on domain)
├── infrastructure/  # In-memory repositories
└── interface/       # Adapters, DTOs, HTTP handlers, router
    ├── http/
    │   ├── dto/
    │   ├── handler/
    │   ├── middleware/
    │   └── router/
    ├── brand/
    ├── category/
    └── ... (one adapter per domain)
cmd/
└── main.go          # Wiring + entry point
tests/
├── entity/          # Domain entity tests
├── application/     # Use case tests
└── interface/       # Adapter + HTTP handler tests
```

## Domains

| Domain | Description | Status |
|---|---|---|
| Store | Seller storefront management | ✅ |
| Inventory | Product inventory, sale models, promotions | ✅ |
| Product | Product catalog with advanced features | ✅ |
| Brand | Product brand management | ✅ |
| Category | Product category management | ✅ |
| Promotion | Discount/promotion management | ✅ |
| Reference Price | Market reference pricing | ✅ |
| Sales Commission | Commission rate configuration | ✅ |
| Store Allowed Category | Category access control per store | ✅ |
| Store–Warehouse Link | Warehouse assignment to stores | ✅ |
| Warehouse | Physical warehouse management | ✅ |
| Product Image | Product gallery (multiple images, ordering) | ✅ |
| Product Type | Product variants (size, color, etc.) | ✅ |
| Product Attribute | Custom key-value fields per product | ✅ |
| Price History | Price change audit trail | ✅ |
| Product Bundle | Bundles, upsells, and cross-sells | ✅ |

## Features

- **Full CRUD** for Store, Inventory, Product, Brand, Category, and 11 other domains
- **Product listing** with filters (owner, status, category, brand, text search, pagination)
- **Seller self-service** — "My Products" endpoint
- **Product SEO** — slug generation, meta title, meta description
- **Admin controls** — enable/disable products
- **Product images** — multiple files with sort ordering
- **Product types** — size, color, and other variant dimensions
- **Product attributes** — custom key/value fields
- **Price history** — track price changes over time
- **Product bundles** — bundle, upsell, and cross-sell relations
- **Inventory management** — stock levels, sale models, reservations, low-stock alerts
- **Promotions** — activate/deactivate promotions on inventory items
- **Pagination** on all list endpoints

## Getting Started

```bash
go build ./...
go test ./tests/...
go run ./cmd/
```

The server starts on `:8080` with all endpoints registered.

## API

See [OpenAPI spec](internal/interface/http/handler/openapi.yaml) for the full API reference.

## Test Coverage

69 test files across 35 test suites — every entity, use case, adapter, and HTTP handler is tested.

```
go test ./tests/... -v
```
