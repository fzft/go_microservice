package data

import (
	"context"
	"fmt"
	protos "github.com/fzft/go_microservice/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64), nil}
	go pb.handleUpdates()
	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("Unable to subscribe for rates", "error", err)
	}

	p.client = sub

	for {
		rr, err := sub.Recv()
		if grpcError := rr.GetError(); grpcError != nil {
			p.log.Error("Error subscribing for rates", "error", grpcError.GetMessage())
			continue
		}

		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("Received updated rate from server", "dest", resp.GetDestination().String())
			if err != nil {
				p.log.Error("Error receiving message", "error", err)
				return
			}
			p.rates[resp.GetDestination().String()] = resp.Rate
		}

	}

}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}
	return pr, nil
}

func GetProducts() Products {
	return productList
}

func (p *ProductsDB) AddProduct(pr Product) {
	pr.ID = getNextID()
	productList = append(productList, &pr)
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(pr.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &pr

	return nil
}

// DeleteProduct deletes a product from the database
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}
	np := *productList[i]
	np.Price = np.Price * rate
	return &np, nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	//if cached return
	//if r, ok := p.rates[destination]; ok {
	//	return r, nil
	//}

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}
	// get initial rate
	resp, err := p.currency.GetRate(context.Background(), rr)
	if err != nil {
		if grpcError, ok := status.FromError(err); ok {
			if grpcError.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to retreive exchange rate from currency service: %s", grpcError.Message())
			}
		}
		return -1, err
	}
	p.rates[destination] = resp.Rate // update cache

	// subscribe for updates
	p.client.Send(rr)

	p.log.Info("resp.Rate", resp.Rate)
	return resp.Rate, err
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
