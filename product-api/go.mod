module github.com/fzft/go_microservice/product-api

go 1.14

require (
	github.com/fzft/go_microservice/currency v0.0.0-20200809153505-122bac1dd353
	github.com/go-openapi/runtime v0.19.20
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/leodido/go-urn v1.2.0 // indirect
	google.golang.org/grpc v1.31.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/fzft/go_microservice/currency => ../currency