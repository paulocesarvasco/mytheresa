package categoriesapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/mytheresa/go-hiring-challenge/internal/api/middlewares"
	"github.com/mytheresa/go-hiring-challenge/internal/payloads"
)

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.With(middlewares.ValidateJSON[payloads.CreateCategoryRequest](h.log)).
		Post("/", h.CreateCategory)

	r.With(middlewares.ParseQueryParameters(h.log)).
		Get("/", h.GetCategories)

	return r
}
