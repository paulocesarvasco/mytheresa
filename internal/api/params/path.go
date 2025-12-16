package params

import "context"

type PathParams struct {
	Code string
}

type pathParamsKey struct{}

func WithPathParams(ctx context.Context, p PathParams) context.Context {
	return context.WithValue(ctx, pathParamsKey{}, p)
}

func PathParamsFromContext(ctx context.Context) (PathParams, bool) {
	p, ok := ctx.Value(pathParamsKey{}).(PathParams)
	return p, ok
}
