package params

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
)

type bodyKey[T any] struct{}

func WithBody[T any](ctx context.Context, v []T) context.Context {
	return context.WithValue(ctx, bodyKey[T]{}, v)
}

func BodyFromRequest[T any](r *http.Request) ([]T, bool) {
	v, ok := r.Context().Value(bodyKey[T]{}).([]T)
	return v, ok
}

type ListOrObject[T any] struct {
	Items []T
}

func (l *ListOrObject[T]) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return errorsapi.ErrEmptyJSONBody
	}

	switch b[0] {
	case '{':
		var single T
		if err := json.Unmarshal(b, &single); err != nil {
			return err
		}
		l.Items = []T{single}
		return nil
	case '[':
		var list []T
		if err := json.Unmarshal(b, &list); err != nil {
			return err
		}
		l.Items = list
		return nil
	default:
		return errorsapi.ErrUnexpectedBodyFormat
	}
}
