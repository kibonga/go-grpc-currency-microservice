package data

import (
	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

//func NewRates(l hclog.Logger) (*ExchangeRates, error) {
//	err := &ExchangeRates{l, map[string]float64{}}
//
//}

//func (er *ExchangeRates) getRates() error {
//	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
//
//	if err != nil {
//		return nil
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("Status code not valid. StatusCode=%d", resp.StatusCode)
//	}
//
//	defer resp.Body.Close()
//
//}
