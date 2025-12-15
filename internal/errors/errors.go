package errorsapi

import "errors"

var (
	ErrCatalogInvalidLimit    = errors.New("invalid limit parameter")
	ErrCatalogInvalidOffset   = errors.New("invalid offset parameter")
	ErrCatalogInvalidMaxPrice = errors.New("invalid max_price parameter")

	ErrInvalidProductCode = errors.New("invalid product code")

	ErrRepositoryCountProducts = errors.New("failed to count products")
	ErrRepositoryFetchProducts = errors.New("failed to fetch products")
)
