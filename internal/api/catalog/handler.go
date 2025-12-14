package catalogapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
)

type Service interface {
	ListProducts(ctx context.Context, limit int, offset int) (catalog.ProductPage, error)
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
	offset := 0
	if v := queryParameters.Get("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 1 {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = parsed
	}
	if v := queryParameters.Get("offset"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}
		offset = parsed
	}

	products, err := h.service.ListProducts(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody := &bytes.Buffer{}
	if err := json.NewEncoder(responseBody).Encode(products); err != nil {
		// TODO: improve error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody.Bytes())

}
