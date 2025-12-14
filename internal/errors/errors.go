package errorsapi

import "errors"

var (
	ErrCatalogInvalidLimit    = errors.New("invalid limit parameter")
	ErrCatalogInvalidOffset   = errors.New("invalid offset parameter")
	ErrCatalogInvalidMaxPrice = errors.New("invalid max_price parameter")

	ErrRepositoryCountProducts = errors.New("failed to count products")
	ErrRepositoryFetchProducts = errors.New("failed to fetch products")
)
