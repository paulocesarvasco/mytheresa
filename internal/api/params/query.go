package params

import (
	"context"

	"github.com/shopspring/decimal"
)

const (
	DefaultLimit  = 10
	DefaultOffset = 0
	MaxLimit      = 100
)

type QueryParams struct {
	Limit        int
	Offset       int
	CategoryCode string
	MaxPrice     *decimal.Decimal
}

type queryParamsKey struct{}

func WithQueryParams(ctx context.Context, p QueryParams) context.Context {
	return context.WithValue(ctx, queryParamsKey{}, p)
}

func QueryParamsFromContext(ctx context.Context) QueryParams {
	if v, ok := ctx.Value(queryParamsKey{}).(QueryParams); ok {
		return v
	}
	return QueryParams{Limit: DefaultLimit, Offset: DefaultOffset}
}
