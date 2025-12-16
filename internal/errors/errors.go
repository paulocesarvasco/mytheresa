package errorsapi

import "errors"

var (
	ErrInvalidLimitParam    = errors.New("invalid limit parameter")
	ErrInvalidOffsetParam   = errors.New("invalid offset parameter")
	ErrInvalidMaxPriceParam = errors.New("invalid max_price parameter")
	ErrInvalidJSONBody      = errors.New("invalid json body")
	ErrInvalidRequestSchema = errors.New("invalid request schema")
	ErrInvalidContentType   = errors.New("content type must be application/json")

	ErrMissingRequestParams = errors.New("request parameters not found in context")

	ErrInvalidProductCode = errors.New("invalid product code")

	ErrRepositoryCountProducts   = errors.New("failed to count products")
	ErrRepositoryFetchProducts   = errors.New("failed to fetch products")
	ErrRepositoryCountCategories = errors.New("failed to count categories")
	ErrRepositoryFetchCategories = errors.New("failed to fetch categories")

	ErrProductNotFound        = errors.New("product not found")
	ErrRepositoryFetchProduct = errors.New("failed to fetch product")
)
