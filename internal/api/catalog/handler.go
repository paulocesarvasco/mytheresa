package catalogapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/shopspring/decimal"
)

type Service interface {
	ListProducts(ctx context.Context, limit, offset int, categoryCode string, maxPrice *decimal.Decimal) (catalog.ProductPage, error)
}

type Handler struct {
	service Service
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	queryParameters := r.URL.Query()

	limit := 10
	if v := queryParameters.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 {
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
			api.ErrorResponse(w, r, http.StatusBadRequest, errorsapi.ErrCatalogInvalidMaxPrice.Error())
			return
		}
		maxPrice = &parsed
	}

	products, err := h.service.ListProducts(r.Context(), limit, offset, categoryCode, maxPrice)
	if err != nil {
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, products)
}
