package catalogapi

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.With(ValidateCatalogQuery(h.log)).
		Get("/", h.GetProducts)

	r.With(ValidateProductCode(h.log)).
		Get("/{code}", h.GetDetailProduct)

	return r
}
