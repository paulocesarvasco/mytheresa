package catalogapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/mytheresa/go-hiring-challenge/internal/api/middlewares"
)

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.With(middlewares.ParseQueryParameters(h.log)).
		With(middlewares.ParseMaxPrice(h.log)).
		Get("/", h.GetProducts)

	r.With(middlewares.ValidateProductCode(h.log)).
		Get("/{code}", h.GetDetailProduct)

	return r
}
