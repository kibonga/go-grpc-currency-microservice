package main

import (
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
	log := hclog.Default()

	grpcServer := grpc.NewServer()
	currencyServer := server.NewCurrency(log)

	pb.RegisterCurrencyServer(grpcServer, currencyServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	fmt.Println("Listenting...")
	grpcServer.Serve(listener)
}
