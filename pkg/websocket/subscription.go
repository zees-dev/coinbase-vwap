package websocket

import (
	"encoding/json"
)

// https://docs.cloud.coinbase.com/exchange/docs/overview
const DefaultCoinbaseURL = "wss://ws-feed.exchange.coinbase.com"

type subscriptionType string

const (
	subscriptions subscriptionType = "subscriptions"
	subscribe     subscriptionType = "subscribe"
	unsubscribe   subscriptionType = "unsubscribe"
)

type supportedChannel string

const (
	matches   supportedChannel = "matches"
	heartbeat supportedChannel = "heartbeat"
	status    supportedChannel = "status"
	ticker    supportedChannel = "ticker"
	level2    supportedChannel = "level2"
	// TODO: add more channels to match coinbase API
)

type channel struct {
	Name       supportedChannel `json:"name"`
	ProductIDs []string         `json:"product_ids"`
}

// subscription represents a request to subscribe or unsubscribe to a channel
// example subscribe: { "type": "subscribe", "channels": [{ "name": "matches", "product_ids": ["BTC-USD"] }] }
// example unsubscribe: { "type": "unsubscribe", "channels": [{ "name": "matches", "product_ids": ["BTC-USD"] }] }
type subscription struct {
	Type     subscriptionType `json:"type"`
	Channels []channel        `json:"channels"`
}

func (s subscription) Byte() []byte {
	b, _ := json.Marshal(s)
	return b
}

// message is the channel response/message from a channal subscription request
// example:
// {"type":"match","trade_id":229219170,"maker_order_id":"b24e6960-0090-4dc6-b626-46f38705ddf7","taker_order_id":"27cd034c-b3d3-4e45-818a-11dd09b71a59","side":"sell","size":"0.0016","price":"61780.37","product_id":"BTC-USD","sequence":30655769767,"time":"2021-10-30T10:48:58.114400Z"}
// note: unused fields are commented out for performance reasons (no need for unmarshalling)
type message struct {
	Type supportedChannel `json:"type"`
	// TradeID      int64            `json:"trade_id"`
	// MakerOrderID string           `json:"maker_order_id"`
	// TakerOrderID string           `json:"taker_order_id"`
	// Side         string           `json:"side"`
	Size      string `json:"size"`
	Price     string `json:"price"`
	ProductID string `json:"product_id"`
	// Sequence     int64            `json:"sequence"`
	// Time         string           `json:"time"`
}

func toMessage(msgB []byte) (message, error) {
	var msg message
	err := json.Unmarshal(msgB, &msg)
	return msg, err
}
