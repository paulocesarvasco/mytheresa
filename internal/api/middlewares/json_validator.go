package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/api/params"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

func ValidateJSON[T any](log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ct := r.Header.Get("Content-Type")
			if !(ct == "application/json") {
				log.Warn(r.Context(), "invalid content type", "content_type", ct)
				api.ErrorResponse(w, r, http.StatusUnsupportedMediaType, errorsapi.ErrInvalidContentType.Error())
				return
			}

			var payload T

			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()

			if err := dec.Decode(&payload); err != nil {
				log.Warn(r.Context(), "invalid json body", "err", err)
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidJSONBody.Error())
				return
			}

			if dec.More() {
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidJSONBody.Error())
				return
			}

			v := validator.New()
			if err := v.Struct(payload); err != nil {
				log.Warn(r.Context(), "schema validation failed", "err", err)
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidRequestSchema.Error())
				return
			}

			ctx := params.WithBody(r.Context(), payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
