package models

type KrakenWebSocketResponse struct {
	Channel string `json:"channel"`
	Type    string `json:"type"`
	Data    []struct {
		Symbol string  `json:"symbol"`
		Last   float64 `json:"last"`
	} `json:"data"`
}

type SubscriptionRequest struct {
	Method             string `json:"method"`
	SubscriptionParams `json:"params"`
}

type SubscriptionParams struct {
	Channel string   `json:"channel"`
	Symbol  []string `json:"symbol"`
}

type KrakenTickerResponse struct {
	Error  []string                `json:"error"`
	Result map[string]KrakenTicker `json:"result"`
}

type KrakenTicker struct {
	C []string `json:"c"`
}
