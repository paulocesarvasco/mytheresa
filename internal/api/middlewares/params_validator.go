package middlewares

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/api/params"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

func ParseQueryParameters(log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			// start from defaults (or previously set params)
			p := params.QueryParamsFromContext(r.Context())

			if v := q.Get("limit"); v != "" {
				parsed, err := strconv.Atoi(v)
				if err != nil || parsed < 1 {
					log.Warn(r.Context(), "invalid query parameter",
						"param", "limit", "value", v, "err", err,
					)
					api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidLimitParam.Error())
					return
				}
				if parsed > params.MaxLimit {
					log.Warn(r.Context(),
						"query parameter \"limit\" exceeds maximum allowed value",
						"limit", parsed,
						"max", params.MaxLimit,
					)
					parsed = params.MaxLimit
				}
				p.Limit = parsed
			}

			if v := q.Get("offset"); v != "" {
				parsed, err := strconv.Atoi(v)
				if err != nil || parsed < 0 {
					log.Warn(r.Context(), "invalid query parameter",
						"param", "offset", "value", v, "err", err,
					)
					api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidOffsetParam.Error())
					return
				}
				p.Offset = parsed
			}

			if v := q.Get("category_code"); v != "" {
				// TODO: validate format
				p.CategoryCode = v
			}

			ctx := params.WithQueryParams(r.Context(), p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ParseMaxPrice(log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.URL.Query().Get("max_price")
			if v == "" {
				next.ServeHTTP(w, r)
				return
			}

			parsed, err := decimal.NewFromString(v)
			if err != nil || !parsed.GreaterThan(decimal.Zero) {
				log.Warn(r.Context(), "invalid query parameter",
					"param", "max_price", "value", v, "err", err,
				)
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidMaxPriceParam.Error())
				return
			}

			p := params.QueryParamsFromContext(r.Context())
			p.MaxPrice = &parsed

			ctx := params.WithQueryParams(r.Context(), p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var productCodeRe = regexp.MustCompile(`^PROD\d{3}$`)

func ValidateProductCode(log logs.ApiLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := chi.URLParam(r, "code")

			if !productCodeRe.MatchString(code) {
				log.Warn(
					r.Context(),
					"product code does not match required pattern",
					"expected_pattern", "PROD###",
					"code", code,
				)
				api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidProductCode.Error())
				return
			}

			ctx := params.WithPathParams(r.Context(), params.PathParams{Code: code})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
