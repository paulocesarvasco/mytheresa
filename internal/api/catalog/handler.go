package catalogapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/shopspring/decimal"
)

type Service interface {
	ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) (catalog.ProductPage, error)
	DetailProduct(ctx context.Context, code string) (catalog.ProductView, error)
}

type Handler struct {
	service Service
	log     *logs.CustomLogger
}

func New(s Service) *Handler {
	logger := logs.NewLogger()
	return &Handler{
		service: s,
		log:     logger,
	}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	queryParameters := r.URL.Query()

	limit := 10
	if v := queryParameters.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 {
			h.log.Warn(r.Context(), "invalid limit parameter", "error", err, "limit", v)
			api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidLimit.Error())
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = parsed
	}

	offset := 0
	if v := queryParameters.Get("offset"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 0 {
			h.log.Warn(r.Context(), "invalid offset parameter", "error", err, "offset", v)
			api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidOffset.Error())
			return
		}
		offset = parsed
	}

	categoryCode := ""
	if v := queryParameters.Get("category_code"); v != "" {
		// TODO: validate input
		categoryCode = v
	}

	var maxPrice *decimal.Decimal
	if v := queryParameters.Get("max_price"); v != "" {
		parsed, err := decimal.NewFromString(v)
		if err != nil || !parsed.GreaterThan(decimal.Zero) {
			h.log.Warn(r.Context(), "invalid max_price parameter", "error", err, "max_price", v)
			api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidMaxPrice.Error())
			return
		}
		maxPrice = &parsed
	}

	products, err := h.service.ListProducts(r.Context(), limit, offset, categoryCode, maxPrice)
	if err != nil {
		h.log.Error(r.Context(), "list products failed", "err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, products)
}

func (h *Handler) GetDetailProduct(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if !isValidProductCode(code) {
		h.log.Warn(r.Context(), "product code does not match required pattern", "expected_pattern", "PROD###", "code", code)
		api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrInvalidProductCode.Error())
		return
	}

	details, err := h.service.DetailProduct(r.Context(), code)
	if err != nil {
		if errors.Is(err, errorsapi.ErrProductNotFound) {
			api.ErrorResponse(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.log.Error(r.Context(), "retrieve product detail failed", "err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, details)
}
