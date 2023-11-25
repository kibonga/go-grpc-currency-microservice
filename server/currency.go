package server

import (
	data "CurrencyMicroservice/data"
	pb "CurrencyMicroservice/protos/currency"
	"context"
	"github.com/hashicorp/go-hclog"
	"io"
	"time"
)

type CurrencyServerApi struct {
	log hclog.Logger
	pb.UnimplementedCurrencyServer
	rates *data.ExchangeRates
}

func NewCurrencyServerApi(l hclog.Logger) *CurrencyServerApi {
	return &CurrencyServerApi{log: l}
}

func (c *CurrencyServerApi) GetRate(ctx context.Context, rateReq *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Getting the Rate", rateReq)

	return &pb.RateResponse{Rate: 0.5}, nil
}

func (c *CurrencyServerApi) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {

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
