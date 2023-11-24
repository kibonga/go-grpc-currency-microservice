package server

import (
	pb "CurrencyMicroservice/protos/currency"
	"context"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Currency struct {
	log hclog.Logger
	pb.UnimplementedCurrencyServer
}

type CurrencySubscribeRatesServer struct {
	log hclog.Logger
	grpc.ServerStream
}

func (c *Currency) GetRate(ctx context.Context, rateReq *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Getting the Rate", rateReq)

	return &pb.RateResponse{Rate: 0.5}, nil
}

func (c *Currency) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {

	// Client streaming request (blocking method, use goroutines)
	go func() {
		for {
			req, err := src.Recv()

			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}

			if err != nil {
				c.log.Error("Unable to read from client")
				break
			}

			c.log.Info("Handle client request", "Base", req.GetBase(), "Destination", req.GetDestination())
		}
	}()

	// Server streaming response
	for {
		msg := &pb.RateResponse{Rate: 12.2}

		err := src.Send(msg)
		if err != nil {
			return err
		}

		time.Sleep(3 * time.Second)
	}
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}
