package server

import (
	"context"
	"fmt"
	"github.com/fzft/go_microservice/currency/data"
	protos "github.com/fzft/go_microservice/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type Currency struct {
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got updated rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
				err = k.Send(&protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
					},
				})
				if err != nil {
					c.log.Error("Unable to send update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())

				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	if rr.Base == rr.Destination {
		err := status.Errorf(codes.InvalidArgument,
			"Base currency %s can not be the same as the destination %s currency",
			rr.GetBase().String(), rr.GetDestination().String(),
		)
		return nil, err
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Base: rr.GetBase(), Destination: rr.GetDestination(), Rate: rate}, nil
}

// SubscribeRates implments the gRPC bidirection streaming method for the server
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	// handle client messages
	for {
		rr, err := src.Recv() // Recv is a blocking method which returns on client data
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			break
		}

		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		// check that subscription does not exist
		var validationError *status.Status
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {
				c.log.Error("Subscription already active", "base", rr.Base.String(), "dest", rr.Destination.String())
				// subscription exists return errors
				validationError = status.New(codes.AlreadyExists, "Unable to subscribe for currency as subscription already exists")
				validationError, err = validationError.WithDetails(rr)
				if err != nil {
					c.log.Error("Unable to add metadata to error", "error", err)
					break
				}
				break
			}
		}
		fmt.Println(validationError)

		if validationError != nil {
			c.log.Error("send validationError", "error", err)
			src.Send(&protos.StreamingRateResponse{Message: &protos.StreamingRateResponse_Error{
				Error: validationError.Proto(),
			}})
			continue
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
