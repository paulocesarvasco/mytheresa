package catalogapi

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
)

type Handler struct {
	service *catalog.Service
}

func New(s *catalog.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
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
