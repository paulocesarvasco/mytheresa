package errorsapi

import "errors"

var (
	ErrCatalogInvalidLimit        = errors.New("invalid limit parameter")
	ErrCatalogInvalidOffset       = errors.New("invalid offset parameter")
	ErrCatalogInvalidMaxPrice     = errors.New("invalid max_price parameter")
	ErrCatalogInvalidContextState = errors.New("invalid request context state")

	ErrInvalidProductCode = errors.New("invalid product code")

	ErrRepositoryCountProducts   = errors.New("failed to count products")
	ErrRepositoryFetchProducts   = errors.New("failed to fetch products")
	ErrRepositoryCountCategories = errors.New("failed to count categories")
	ErrRepositoryFetchCategories = errors.New("failed to fetch categories")

	ErrProductNotFound        = errors.New("product not found")
	ErrRepositoryFetchProduct = errors.New("failed to fetch product")
)
