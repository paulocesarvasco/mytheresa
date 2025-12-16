package catalogapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/api/params"
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
	log     logs.ApiLogger
}

func New(s Service) *Handler {
	return &Handler{
		service: s,
		log:     logs.Logger(),
	}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	p := params.QueryParamsFromContext(r.Context())

	products, err := h.service.ListProducts(r.Context(), p.Limit, p.Offset, p.CategoryCode, p.MaxPrice)
	if err != nil {
		h.log.Error(r.Context(), "get products list failed",
			"err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, products)
}

func (h *Handler) GetDetailProduct(w http.ResponseWriter, r *http.Request) {
	p, ok := params.PathParamsFromContext(r.Context())
	if !ok {
		api.ErrorResponse(w, r, http.StatusInternalServerError, errorsapi.ErrCatalogInvalidContextState.Error())
		return
	}

	details, err := h.service.DetailProduct(r.Context(), p.Code)
	if err != nil {
		if errors.Is(err, errorsapi.ErrProductNotFound) {
			api.ErrorResponse(w, r, http.StatusNotFound, err.Error())
			return
		}
		h.log.Error(r.Context(), "retrieve product detail failed",
			"err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, details)
}
