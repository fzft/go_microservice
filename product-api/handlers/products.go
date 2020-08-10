package handlers

import (
	"github.com/fzft/go_microservice/product-api/data"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

// Products handler for getting and updating products
type Products struct {
	l  hclog.Logger
	v  *data.Validation
	productDB *data.ProductsDB
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct {
}
