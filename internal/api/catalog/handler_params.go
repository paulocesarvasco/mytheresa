package catalogapi

import (
	"net/http"

	"github.com/shopspring/decimal"
)

type queryParams struct {
	Limit        int
	Offset       int
	CategoryCode string
	MaxPrice     *decimal.Decimal
}

type queryParamsKey struct{}

func getQueryParams(r *http.Request) queryParams {
	return r.Context().Value(queryParamsKey{}).(queryParams)

}

type pathParams struct {
	Code string
}

type pathParamsKey struct{}

func getPathParams(r *http.Request) pathParams {
	return r.Context().Value(pathParamsKey{}).(pathParams)
}
