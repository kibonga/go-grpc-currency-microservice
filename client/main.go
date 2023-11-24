package main

import (
	pb "CurrencyMicroservice/protos/currency"
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"os"
)

func main() {

	// Sets grpc channel options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	// Creates a grpc channel
	channel, err := grpc.Dial("localhost:9092", opts...)

	if err != nil {
		hclog.Default().Error("Failed to grpc.Dial. Error=", err)
		os.Exit(1)
	}

	// Creates a new grpc service client
	currencyClient := pb.NewCurrencyClient(channel)

	resp, err := currencyClient.GetRate(context.Background(), &pb.RateRequest{Base: pb.Currencies_USD, Destination: pb.Currencies_AUD})

	if err != nil {
		hclog.Default().Error("Failed to .GetRate. Error=", err)
		os.Exit(1)
	}

	fmt.Println("Response Rate = ", resp.Rate)

	defer channel.Close()
}
