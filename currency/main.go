package main

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	protos "github.com/fzft/go_microservice/currency/protos/currency"
	"github.com/fzft/go_microservice/currency/server"
)

func main() {
	log := hclog.Default()
	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	cs := server.NewCurrency(log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gs.Serve(l)

}
