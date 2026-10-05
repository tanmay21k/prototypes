// Package products contains HTTP handlers and service logic for product operations.
package products

import (
	"log"
	"net/http"

	jsonUtil "github.com/tanmay21k/protypes/sqlcDB/internal/json"
)

type handler struct {
	service Service
}

// NewHandler creates a product handler bound to the provided service.
func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonUtil.Write(w, http.StatusOK, products)
}
