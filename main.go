package main

import (
	"context"
	"log"
	"os"

	"github.com/zees-dev/coinbase-vwap/pkg/websocket"
)

func main() {
	// context for the websocket connection
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	vwapSocket := websocket.NewVWAPSocket(os.Stdout)

	// initiate non-blocking client-server websocket communication
	pairs, windowSize := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}, 200
	go vwapSocket.ComputeVWAP(ctx, windowSize, pairs...)

	wsClient, err := websocket.NewWSClient(websocket.DefaultCoinbaseURL, vwapSocket)
	if err != nil {
		log.Fatal(err)
	}

	err = wsClient.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
