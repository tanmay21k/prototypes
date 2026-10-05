package products

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/tanmay21k/protypes/sqlcDB/internal/adapters/postgresql/sqlc"
	jsonUtil "github.com/tanmay21k/protypes/sqlcDB/internal/json"
)

// Service defines the application contract for product operations.
type Service interface {
	ListProducts(context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

// NewService wires a sqlc-backed repo into the product service layer.
func NewService(repo *repo.Querier) Service {
	return &svc{repo: *repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductById(ctx, id)
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindProductById(r.Context(), productID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonUtil.Write(w, http.StatusOK, product)
}
