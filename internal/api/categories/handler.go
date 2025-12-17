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
	ListCategories(ctx context.Context, limit, offset int, categoryCode string) ([]categories.Category, int64, error)
	CreateCategory(ctx context.Context, code string, name string) error
	CreateCategories(ctx context.Context, categories []categories.CreateCategoryInput) error
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

	res, total, err := h.service.ListCategories(r.Context(), p.Limit, p.Offset, p.CategoryCode)
	if err != nil {
		h.log.Error(r.Context(), "list categories failed",
			"err", err)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	categories := make([]CategoryView, len(res))
	for i, c := range res {
		categories[i] = CategoryView{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, r, CategoryPage{Categories: categories, Total: total})
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	req, ok := params.BodyFromRequest[payloads.CreateCategoryRequest](r)
	if !ok {
		api.ErrorResponse(w, r, http.StatusInternalServerError, "invalid request state")
		return
	}

	var err error
	if len(req) == 1 {
		err = h.service.CreateCategory(r.Context(), req[0].Code, req[0].Name)
	} else {
		inputs := make([]categories.CreateCategoryInput, len(req))
		for i, r := range req {
			inputs[i] = categories.CreateCategoryInput{
				Code: r.Code,
				Name: r.Name,
			}
		}
		err = h.service.CreateCategories(r.Context(), inputs)
	}

	if err != nil {
		if errors.Is(err, errorsapi.ErrRepositoryCategoryAlreadyExists) {
			api.ErrorResponse(w, r, http.StatusConflict, err.Error())
			return
		}
		h.log.Error(r.Context(), "create category failed", "err", err, "code", req)
		api.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponseWithStatus(w, r, http.StatusCreated, nil)
}
