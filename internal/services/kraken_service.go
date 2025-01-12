package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"test-assigment-ltp/internal/cache"
	"test-assigment-ltp/internal/models"
)

type KrakenServiceInterface interface {
	GetLastTradedPrice(pair string) (float64, error)
}

type KrakenService struct {
	httpClient *http.Client
	wsConn     *websocket.Conn
	wsMutex    sync.Mutex
	redisCache *cache.RedisCache
}

func NewKrakenService(redisCache *cache.RedisCache) *KrakenService {
	service := &KrakenService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		redisCache: redisCache,
	}

	go service.connectWebSocket()

	return service
}

// connectWebSocket Connect and listen for WebSocket messages
func (ks *KrakenService) connectWebSocket() {
	ks.wsMutex.Lock()
	defer ks.wsMutex.Unlock()

	var err error
	ks.wsConn, _, err = websocket.DefaultDialer.Dial("wss://ws.kraken.com/v2", nil)
	if err != nil {
		logrus.Errorf("Failed to connect to Kraken WebSocket: %v", err)
		return
	}

	go ks.listenWebSocket()
}

// subscribeToPair Subscribe to retrive data for pair
func (ks *KrakenService) subscribeToPair(pair string) {
	request := models.SubscriptionRequest{
		Method: "subscribe",
		SubscriptionParams: models.SubscriptionParams{
			Channel: "ticker",
			Symbol:  []string{pair},
		},
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		logrus.Errorf("failed to marshal subscription request: %v", err)
	}

	if err := ks.wsConn.WriteMessage(websocket.TextMessage, requestBytes); err != nil {
		logrus.Errorf("failed to send subscription request: %v", err)
	}

	logrus.Infof("Subscribed to ticker for symbol: %s", pair)
}

// listenWebSocket listen for WebSocket messages and store it to Redis
func (ks *KrakenService) listenWebSocket() {
	for {
		_, message, err := ks.wsConn.ReadMessage()
		if err != nil {
			logrus.Errorf("WebSocket read error: %v", err)
			ks.reconnectWebSocket()
			return
		}

		var data models.KrakenWebSocketResponse
		if err := json.Unmarshal(message, &data); err != nil {
			logrus.Errorf("Failed to unmarshal WebSocket message: %v", err)
			continue
		}

		// Handle WebSocket events
		switch data.Channel {
		case "heartbeat":
			// Ignore heartbeat messages
			continue
		case "ticker":
			// Process the ticker update
			for _, data := range data.Data {
				lastPrice := data.Last
				pair := data.Symbol

				if err := ks.redisCache.Set(pair, lastPrice, 1*time.Minute); err != nil {
					logrus.Errorf("Failed to update Redis cache for pair %s: %v", pair, err)
				} else {
					logrus.Infof("Updated last price for pair %s: %.2f", pair, lastPrice)
				}
			}
		default:
			logrus.Warnf("Unknown WebSocket event: %s", data.Channel)
		}
	}
}

// reconnectWebSocket() reccontion WebSocket handler
func (ks *KrakenService) reconnectWebSocket() {
	time.Sleep(5 * time.Second)
	ks.connectWebSocket()
}

// GetLastTradedPrice fetch data from Redis or use HTTP API if cache is empty
func (ks *KrakenService) GetLastTradedPrice(pair string) (float64, error) {
	amount, err := ks.redisCache.Get(pair)
	if err == nil {
		return amount, nil
	}

	amount, err = ks.fetchLastTradedPrice(pair)
	if err != nil {
		return 0, err
	}

	go ks.subscribeToPair(pair)

	if err := ks.redisCache.Set(pair, amount, 1*time.Minute); err != nil {
		logrus.Errorf("Failed to update Redis cache for pair %s: %v", pair, err)
	}

	return amount, nil
}

// fetchLastTradedPrice HTTP API method if cache is empty
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
	fmt.Sscanf(ticker.C[0], "%f", &amount)

	return amount, nil
}
