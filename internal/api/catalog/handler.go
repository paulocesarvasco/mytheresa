package catalogapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/internal/catalog"
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
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
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
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
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
			http.Error(w, "invalid max_price parameter", http.StatusBadRequest)
			return
		}
		maxPrice = &parsed
	}

	products, err := h.service.ListProducts(r.Context(), limit, offset, categoryCode, maxPrice)
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
