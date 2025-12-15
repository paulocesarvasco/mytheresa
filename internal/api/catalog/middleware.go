package catalogapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

func ValidateCatalogQuery(log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			p := queryParams{
				Limit:  10,
				Offset: 0,
			}

			if v := q.Get("limit"); v != "" {
				parsed, err := strconv.Atoi(v)
				if err != nil || parsed < 1 {
					log.Warn(r.Context(), "invalid query parameter",
						"param", "limit",
						"value", v,
						"err", err)
					api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidLimit.Error())
					return
				}
				if parsed > 100 {
					parsed = 100
				}
				p.Limit = parsed
			}

			if v := q.Get("offset"); v != "" {
				parsed, err := strconv.Atoi(v)
				if err != nil || parsed < 0 {
					log.Warn(r.Context(), "invalid query parameter",
						"param", "offset",
						"value", v,
						"err", err)
					api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidOffset.Error())
					return
				}
				p.Offset = parsed
			}

			if v := q.Get("category_code"); v != "" {
				p.CategoryCode = v
			}

			if v := q.Get("max_price"); v != "" {
				parsed, err := decimal.NewFromString(v)
				if err != nil || !parsed.GreaterThan(decimal.Zero) {
					log.Warn(r.Context(), "invalid query parameter",
						"param", "max_price",
						"value", v,
						"err", err)
					api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidMaxPrice.Error())
					return
				}
				p.MaxPrice = &parsed
			}

			ctx := context.WithValue(r.Context(), queryParamsKey{}, p)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateProductCode(log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "code")

			if !isValidProductCode(code) {
				log.Warn(r.Context(),
					"product code does not match required pattern",
					"expected_pattern", "PROD###",
					"code", code,
				)
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidProductCode.Error())
				return
			}

			p := pathParams{Code: code}
			ctx := context.WithValue(r.Context(), pathParamsKey{}, p)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
