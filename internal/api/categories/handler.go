package categoriesapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/internal/api"
	"github.com/mytheresa/go-hiring-challenge/internal/api/params"
	"github.com/mytheresa/go-hiring-challenge/internal/categories"
	errorsapi "github.com/mytheresa/go-hiring-challenge/internal/errors"
	"github.com/mytheresa/go-hiring-challenge/internal/logs"
	"github.com/mytheresa/go-hiring-challenge/internal/payloads"
)

type Service interface {
	ListCategories(ctx context.Context, limit, offset int, categoryCode string) (categories.CategoryPage, error)
	CreateCategory(ctx context.Context, code string, name string) (categories.CategoryView, error)
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

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	req, ok := params.BodyFromRequest[payloads.CreateCategoryRequest](r)
	if !ok {
		api.ErrorResponse(w, r, http.StatusInternalServerError, "invalid request state")
		return
	}

	created, err := h.service.CreateCategory(r.Context(), req.Code, req.Name)
	if err != nil {
		if errors.Is(err, errorsapi.ErrRepositoryCategoryAlreadyExists) {
			api.ErrorResponse(w, r, http.StatusConflict, err.Error())
			return
		}
		h.log.Error(r.Context(), "create category failed", "err", err, "code", req.Code)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponseWithStatus(w, r, http.StatusCreated, created)
}
