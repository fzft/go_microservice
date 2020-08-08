package handlers

import (
	"context"
	"github.com/fzft/go_microservice/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		p.getProducts(rw, r)
//		return
//	}
//	if r.Method == http.MethodPost {
//		p.addProduct(rw, r)
//		return
//	}
//	// handle update
//	if r.Method == http.MethodPut {
//		// expect id in the uri
//		reg := regexp.MustCompile("/([0-9]+)")
//		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
//		p.l.Println(g)
//		if len(g) != 1 {
//			http.Error(rw, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		if len(g[0]) != 2 {
//			http.Error(rw, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		idString := g[0][1]
//		id, err := strconv.Atoi(idString)
//		if err != nil {
//			http.Error(rw, "Invalid URL", http.StatusBadRequest)
//			return
//		}
//		p.l.Println("got id", id)
//		p.updateProduct(id, rw, r)
//	}
//
//	// catch all
//	rw.WriteHeader(http.StatusMethodNotAllowed)
//}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()
	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)

}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "id not qualified", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Prod: %#v", prod)
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {
}

func (p *Products) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
