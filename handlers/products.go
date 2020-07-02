package handlers

import (
	"go_micorservices/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// NewProducts returns Products struct
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(res, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(res, r)
		return
	}

	if r.Method == http.MethodPut {
		// expect the id in the URI
		rgx := regexp.MustCompile(`/([0-9]+)`)
		g := rgx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(res, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(res, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(res, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, res, r)

		return
	}

	// catch all
	res.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(res http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")

	lp := data.GetProducts()
	err := lp.ToJSON(res)

	if err != nil {
		http.Error(res, "Unable to marshal", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(res http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(res, "Unable to unmarshal payload", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, res http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(res, "Unable to unmarshal payload", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)

	if err != nil {
		http.Error(res, "Product Not found", http.StatusBadRequest)
		return
	}
}
