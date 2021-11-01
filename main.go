package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/zees-dev/coinbase-vwap/pkg/coinbase"
	"github.com/zees-dev/coinbase-vwap/pkg/websocket"
)

func main() {
	coinbaseWSURL := flag.String("coinbase-ws-url", coinbase.DefaultCoinbaseURL, "Coinbase Websockets API endpoint URL")
	tradingPairs := flag.String("trading-pairs", "BTC-USD,ETH-USD,ETH-BTC", "Comma separated list of product_id pairs. e.g. BTC-USD,BTC-GBP,BTC-EUR,ETH-BTC")
	windowSize := flag.Uint("window", 200, "No. of data points included in the sliding window")
	flag.Parse()

	if *coinbaseWSURL == "" {
		log.Fatal("'coinbase-ws-url' is required")
	}

	if *tradingPairs == "" {
		log.Fatal("'trading-pairs' is required")
	}

	if *windowSize == 0 {
		log.Fatal("'window' of non-zero size is required")
	}

	// context for the websocket connection
	ctx, cancel := context.WithCancel(context.TODO())

	// handle OS interrupts
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt
		log.Println("Shuting down...")
		cancel()
	}()

	vwapSocket := coinbase.NewVWAPSocket(os.Stdout)

	// initiate non-blocking client-server websocket communication
	go vwapSocket.ComputeVWAP(ctx, int(*windowSize), strings.Split(*tradingPairs, ",")...)

	wsClient, err := websocket.NewWSClient(*coinbaseWSURL, vwapSocket)
	if err != nil {
		log.Fatal(err)
	}

	err = wsClient.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
