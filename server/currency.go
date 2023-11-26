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
	// This is a sample subscription cache
	// Each SubscribeRatesServer represents single client that has subscribed to some rates
	// Each SubscribeRatesServer uniquely identifies the specific client
	subscriptions map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest
}

func NewCurrencyServerApi(rates *data.ExchangeRates, l hclog.Logger) *CurrencyServerApi {
	serverApi := &CurrencyServerApi{
		log:           l,
		rates:         rates,
		subscriptions: make(map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest, 0),
	}

	// This goroutine is a common pattern for receiving signal
	// It simulates rate monitoring and rate changes
	go serverApi.handleUpdates()

	return serverApi
}

func (serverApi *CurrencyServerApi) GetRate(ctx context.Context, rateReq *pb.RateRequest) (*pb.RateResponse, error) {
	serverApi.log.Info("Getting the Rate", rateReq)

	rate, err := serverApi.rates.GetRate(rateReq.GetBase().String(), rateReq.GetDestination().String())

	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{
		Base:        rateReq.GetBase(),
		Destination: rateReq.GetDestination(),
		Rate:        rate,
	}, nil
}

func (serverApi *CurrencyServerApi) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {

	for {
		rateReq, err := src.Recv()

		// Client has closed connection
		if err == io.EOF {
			serverApi.log.Info("Client has closed connection")
			break
		}

		// Actual streaming error
		if err != nil {
			serverApi.log.Error("Unable to read from client")
			return err
		}

		serverApi.log.Info("Handle client request", "Base", rateReq.GetBase(), "Destination", rateReq.GetDestination())

		clientSubs, ok := serverApi.subscriptions[src]
		if !ok {
			clientSubs = make([]*pb.RateRequest, 0)
		}

		clientSubs = append(clientSubs, rateReq)
		serverApi.subscriptions[src] = clientSubs
	}

	return nil
}

func (serverApi *CurrencyServerApi) handleUpdates() {
	// This is a fn which will be run as goroutine
	// This is a common signal receiving pattern
	updateRates := serverApi.rates.MonitorRates(time.Second * 5)
	for range updateRates {
		serverApi.log.Info("Rates were updated")

		// Loop over clients and their rate request subscriptions
		for client, rateRequests := range serverApi.subscriptions {
			// Loop over rateRequests for each client
			for _, rateReq := range rateRequests {
				rate, err := serverApi.rates.GetRate(rateReq.GetBase().String(), rateReq.GetDestination().String())

				if err != nil {
					serverApi.log.Error("Failed to get update rate.", "base", rateReq.GetBase(), "destination", rateReq.GetDestination())
				}

				newRate := &pb.RateResponse{
					Base:        rateReq.GetBase(),
					Destination: rateReq.GetDestination(),
					Rate:        rate,
				}

				err = client.Send(newRate)
				if err != nil {
					serverApi.log.Error("Failed to send new rate response to client", "newRate", newRate)
				}
			}
		}
	}
}
