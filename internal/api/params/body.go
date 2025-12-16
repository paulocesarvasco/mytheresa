package params

import (
	"context"
	"net/http"
)

type bodyKey[T any] struct{}

func WithBody[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, bodyKey[T]{}, v)
}

func BodyFromRequest[T any](r *http.Request) (T, bool) {
	v, ok := r.Context().Value(bodyKey[T]{}).(T)
	return v, ok
}
