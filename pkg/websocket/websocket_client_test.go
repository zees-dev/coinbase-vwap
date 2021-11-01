package websocket

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zees-dev/coinbase-vwap/pkg/coinbase"
)

type testPubSuber struct {
	pub chan []byte
	sub chan []byte
}

func (t *testPubSuber) Send() chan []byte {
	return t.pub
}

func (t *testPubSuber) Recv() chan []byte {
	return t.sub
}

// Test_wsClient is an integration test for the websocket client connection with coinbase URL.
// It is expected to connect to the websocket and receive a message from the server.
func Test_wsClient(t *testing.T) {
	is := assert.New(t)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	pubSub := &testPubSuber{make(chan []byte), make(chan []byte)}
	wsClient, err := NewWSClient(coinbase.DefaultCoinbaseURL, pubSub)
	if err != nil {
		log.Fatal(err)
	}

	// channel to capture initial subsciption response
	subCh := make(chan coinbase.Subscription, 1)

	// non-blocking subscribe to the websocket
	go func() {
		pubSub.Send() <- []byte(`{ "type": "subscribe", "channels": [{ "name": "matches", "product_ids": ["BTC-USD"] }] }`)

		var sub coinbase.Subscription
		json.Unmarshal(<-pubSub.Recv(), &sub)
		subCh <- sub
		cancel()
	}()

	err = wsClient.Start(ctx)
	is.NoError(err)

	// wait for the subscription and verify the response
	sub := <-subCh
	is.Equal(coinbase.Subscriptions, sub.Type)
	is.NotEmpty(sub.Channels)
}
