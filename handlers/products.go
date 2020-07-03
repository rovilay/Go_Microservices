package handlers

import (
	"context"
	"fmt"
	"go_micorservices/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// NewProducts returns a new Products instance
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts ...
func (p *Products) GetProducts(res http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")

	lp := data.GetProducts()
	err := lp.ToJSON(res)

	if err != nil {
		http.Error(res, "Unable to marshal", http.StatusInternalServerError)
	}
}

// AddProduct ...
func (p *Products) AddProduct(res http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

// UpdateProduct ...
func (p *Products) UpdateProduct(res http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(res, "Invalid Id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT products", id)

	// get the product from ctx and cast it to Product type
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)

	if err != nil {
		http.Error(res, "Product Not found", http.StatusNotFound)
		return
	}
}

// MiddlewareProductValidation ...
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(res, "Unable to unmarshal payload", http.StatusBadRequest)
			return
		}

		// validate product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				res,
				fmt.Sprintf("Invalid product payload: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}
