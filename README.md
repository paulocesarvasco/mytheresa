# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including CRUD operations and database seeding with initial data.


## Project Structure


### `cmd/`
Contains the application entry points.

- `server/main.go`: Main entry point that starts the HTTP REST API.
- `seed/main.go`: Command used to seed the database with initial product data.

### `internal/`
Contains the **core application logic**.
Packages inside `internal/` cannot be imported by external modules, enforcing encapsulation as recommended by Go conventions.


#### `internal/api/`
HTTP layer responsible for routing, request handling, validation, and response formatting.

- **catalog/**: HTTP handlers and routes related to catalog endpoints.
- **categories/**: HTTP handlers, routes, and response views for category endpoints.
- **middlewares/**: Shared HTTP middlewares such as:
  - JSON and parameter validation
  - Request logging
- **params/**: Typed helpers to extract request parameters from body, path, and query.
- **response.go**: Standardized HTTP response helpers.


#### `internal/catalog/`
Catalog orchestration layer.

- Coordinates catalog-related operations
- Maps repository models to service-level views

#### `internal/categories/`
Category domain service layer.

- Implements business logic for categories
- Supports single and batch operations
- Defines service-level domain types and inputs

#### `internal/database/`
Database infrastructure and configuration.

- PostgreSQL connection initialization
- Test utilities (e.g., PostgreSQL test containers)

#### `internal/errors/`
Centralized application and domain error definitions to ensure consistent error handling across all layers.

#### `internal/logs/`
Logging abstractions and configuration.

- Structured logging using `slog`
- Log-level configuration via environment variables
- Shared logging interfaces

#### `internal/payloads/`
HTTP request payload definitions.

- Used exclusively by handlers for decoding and validation
- Explicitly separated from domain and repository models

#### `internal/repository/`
Persistence layer responsible for database access.

- GORM models for categories, products, and variants
- Repository implementations
- Transaction handling
- Repository-level tests and helpers

### `migrations/`
Contains database migration scripts used to manage schema evolution.

### `.env`
Local environment configuration file used during development.


## Architectural Notes

- Handlers remain thin and focused on HTTP concerns.
- Business rules live in the service layers.
- Database access and transactions are isolated in the repository layer.
- Payloads, domain types, and views are explicitly separated.
- Batch operations and transactional consistency are enforced outside the HTTP layer.


## Application Setup

### Requirements
- Go installed locally.
- Docker installed locally.
- Makefile Targets:
  - `make tidy`: Manage and synchronize Go dependencies.
  - `make run`: Start the application locally.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Run the full test suite.
  - `make coverage`: Generate and open a test coverage report.
  - `make docker-up`: Start the application infrastructure using containers.
  - `make docker-up-db`: Start only the database service using containers.
  - `make docker-down`: Stop the Docker infrastructure.

---

Follow up for the assignment here:
[ASSIGNMENT.md](ASSIGNMENT.md)
