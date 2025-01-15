package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"test-assigment-ltp/internal/models"
	"test-assigment-ltp/internal/services"

	"github.com/stretchr/testify/assert"
)

// mock http client

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestFetchLastTradedPrice(t *testing.T) {
	t.Run("Successful response", func(t *testing.T) {
		mockResponse := models.KrakenTickerResponse{
			Error: []string{},
			Result: map[string]models.KrakenTicker{
				"BTC/USD": {
					C: []string{"150000.0"},
				},
			},
		}
		responseBody, _ := json.Marshal(mockResponse)

		client := NewTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "/0/public/Ticker", req.URL.Path)
			assert.Equal(t, "BTC/USD", req.URL.Query().Get("pair"))
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}
		})

		service := services.NewKrakenService(client)

		price, err := service.GetLastTradedPrice("BTC/USD")

		assert.NoError(t, err)
		assert.Equal(t, 150000.0, price)
	})

	t.Run("Unexpected status code", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "/0/public/Ticker", req.URL.Path)
			assert.Equal(t, "BTC/USD", req.URL.Query().Get("pair"))
			return &http.Response{
				StatusCode: 500,
				Header:     make(http.Header),
			}
		})

		service := services.NewKrakenService(client)

		price, err := service.GetLastTradedPrice("BTC/USD")

		assert.Zero(t, price)
		assert.EqualError(t, err, "unexpected status code: 500")
	})
}
