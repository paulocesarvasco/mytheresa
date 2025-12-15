package categoriesapi

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetCategories)

	return r
}
