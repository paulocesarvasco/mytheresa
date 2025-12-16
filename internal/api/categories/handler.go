package categoriesapi

import (
	"context"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/api/params"
	"github.com/mytheresa/go-hiring-challenge/internal/categories"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
)

type Service interface {
	ListCategories(ctx context.Context, limit, offset int, categoryCode string) (categories.CategoryPage, error)
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

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	p := params.QueryParamsFromContext(r.Context())

	products, err := h.service.ListCategories(r.Context(), p.Limit, p.Offset, p.CategoryCode)
	if err != nil {
		h.log.Error(r.Context(), "list categories failed",
			"err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, r, products)
}
