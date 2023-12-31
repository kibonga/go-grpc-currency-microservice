package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Cubes struct {
	Cubes []Cube `xml:"Cube>Cube>Cube"`
}

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	// Get the initial rates
	rates := &ExchangeRates{l, map[string]float64{}}

	err := rates.getRates()

	return rates, err
}

func (er *ExchangeRates) GetRate(base, dest string) (float64, error) {

	baseRate, exist := er.rates[base]
	if !exist {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}

	destRate, exist := er.rates[dest]
	if !exist {
		return 0, fmt.Errorf("Rate not found for currency %d", dest)
	}

	return (destRate / baseRate), nil

}

// Creates a channel used for signaling purposes (struct {} is equivalent to void)
// This fn simulates the changes in the rates (daily rate fluctuation)
func (er *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	channel := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:

				// Adding random diff to the rate to simulate fluctuations in currency rates

				for k, v := range er.rates {
					// Up to 10% difference in rate value
					diff := rand.Float64() / 10
					// Positive or negative change
					dir := rand.Intn(1)

					if dir == 0 {
						diff = 1 - diff
					} else {
						diff = 1 + diff
					}

					// Modify the actual rate
					er.rates[k] = v * diff
				}

				channel <- struct{}{}
			}
		}

	}()

	return channel
}

func (er *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status code not valid. StatusCode=%d", resp.StatusCode)
	}

	defer resp.Body.Close()

	cubes := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(cubes)

	for _, c := range cubes.Cubes {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		er.rates[c.Currency] = r
	}

	er.rates["EUR"] = 1

	return nil
}
