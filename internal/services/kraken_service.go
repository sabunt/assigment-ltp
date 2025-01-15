package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"test-assigment-ltp/internal/models"
)

type KrakenServiceInterface interface {
	GetLastTradedPrice(pair string) (float64, error)
}

type KrakenService struct {
	httpClient *http.Client
}

func NewKrakenService(client *http.Client) *KrakenService {
	service := &KrakenService{
		httpClient: client,
	}
	return service
}

// GetLastTradedPrice fetch data from Kraken API
func (ks *KrakenService) GetLastTradedPrice(pair string) (float64, error) {
	amount, err := ks.fetchLastTradedPrice(pair)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func (ks *KrakenService) fetchLastTradedPrice(pair string) (float64, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=" + pair
	resp, err := ks.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.KrakenTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Error) > 0 {
		return 0, fmt.Errorf("error from Kraken API: %v", result.Error)
	}

	ticker, ok := result.Result[pair]
	if !ok {
		return 0, fmt.Errorf("pair %s not found in Kraken response", pair)
	}

	if len(ticker.C) == 0 {
		return 0, fmt.Errorf("no last trade data for pair %s", pair)
	}

	var amount float64
	if _, err := fmt.Sscanf(ticker.C[0], "%f", &amount); err != nil {
		return 0, fmt.Errorf("failed failed to parse LTP for pair %s: %v", pair, err)
	}

	return amount, nil
}
