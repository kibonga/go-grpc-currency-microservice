package main

import (
	data "CurrencyMicroservice/data"
	pb "CurrencyMicroservice/protos/currency"
	"CurrencyMicroservice/server"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	// Creates default logger instance
	log := hclog.Default()

	rates, err := data.NewRates(log)

	if err != nil {
		os.Exit(1)
	}

	base := "USD"
	dest := "ISK"
	convRate, err := rates.GetRate(base, dest)

	if err != nil {
		fmt.Println("Cannot convert rate for base", base, "and dest", dest, "error=", err)
		os.Exit(1)
	}

	fmt.Println("Base", base, "Dest", dest, "Conversion rate", convRate)

	// Create new empty grpc server
	grpcServer := grpc.NewServer()
	// Creates a server side api for Currency service (contains all server side methods)
	currencyServerApi := server.NewCurrencyServerApi(rates, log)
	// Register service with the grpc server (so when request hits grpc server it will know which service method to call)
	pb.RegisterCurrencyServer(grpcServer, currencyServerApi)
	// Enable grpcurl using reflection api
	reflection.Register(grpcServer)

	// Listens the network over TCP on port 9092
	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	fmt.Println("Listenting...")
	// Start grpc server and accept incoming network connections from the net listener
	grpcServer.Serve(listener)
}
