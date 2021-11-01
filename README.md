# Coinbase Realtime Volume-Weighted-Average-Price Engine

A realtime VWAP (volume-weighted average price) calculation engine using the trading pair data from the [coinbase websocket](https://docs.pro.coinbase.com/#the-matches-channel).

This calculates the [VWAP](https://en.wikipedia.org/wiki/Volume-weighted_average_price) of crypto coin trading pairs as the data from the [match channel](https://docs.cloud.coinbase.com/exchange/docs/channels#match) becomes available (using websockets).

An arbitrary no. of trading pairs could be supported (by specifying them as CLI args); in our case, we default to the following 3:

- `BTC-USD`
- `ETH-USD`
- `ETH-BTC`

## Design

### Key features

- Support for specifying arbitrary no. of trading pairs as CLI args (default to the 3 above)
- The VWAP will only be applicable to the first **200** data points per trading pair (configurable as CLI arg)
- The output will be streamed to a configurable `io.Writer` (`os.Stdout` in our case)
- Generic websocket client implementation - a [PubSuber](pkg/websocket/websocket.go#L4) can read from and write to (to communicate with the server)
  - Look at [tests](pkg/websocket/websocket_client_test.go) for more details
  - The design is extensible since different `PubSuber` implementations could be used to initate different types of subscriptions (i.e. to different coinbase channels, etc.) by communicating with the connected websocket
- Use of the [decimal](https://github.com/shopspring/decimal) for high-precision floating-point arithmatic/calculations (a must in finance)

### Details

The focus of the implementation is to create a generic sliding-window VWAP calculator by consuming price/volume data from a websocket connection.

The implementation consists of 3 core components:

- VWAP calculator - `pkg/vwap`
  - A generic sliding-window VWAP calculator
    - A generic queue from `pkg/queue` provides the _sliding-window_ functionality
  - The calculator retains a running calculation of the current VWAP (as it is updated)
  - The price/volume data is stored as [decimal](https://github.com/shopspring/decimal)
  - If the number of incoming datapoints exceeds the sliding-window (queue) size, the datapoints in the head of the window are dropped as the window moves forward
- Websocket - `pkg/websocket`
  - A generic websocket client setup that wraps around [gorilla/websocket](https://github.com/gorilla/websocket). An implementation of a [PubSuber](pkg/websocket/websocket.go#L4) can be used to send and recieve data from the websocket.
- Coinbase subscriptions (pub/sub) - `pkg/coinbase`
  - The package defines the primitives/structs used by the coinbase subscriptions
  - It also includes an implementation of the `PubSuber` which subscribes to the public [Coinbase websocket url](wss://ws-feed.exchange.coinbase.com)
    - This will subscribe to the matches channel, create a VWAP calculator for each of the trading-pairs of interest (provided as function params), and perform/update the VWAP calculation for the respective trading pairs (as soon as data from websocket is available/recieved)
    - The PubSuber is also responsible for writing data/responses to an `io.Writer` (`STDOUT` in our case).
  - This is essentially a link between VWAP calculator and Websocket components defined above.

The [main.go](main.go) file bootstraps the project by setting up (instantiating) the core components and then running the blocking coinbase websocket connection.

Note: Due to the simplicity of the project, the default std logger is used to debug/display any output.

### Assumptions

- The `side` (buying or selling of trading pair) is ignored; the VWAP calculation only takes into account the volume and price of the trading pair
- The VWAP calculation drops data points which lie outside the sliding window; this may skew the result of the VWAP if the data point(s) dropped have huge volume (relative to all the other data points in the sliding window)

## Run the project

Running the project will initialize a websocket connection to the coinbase websocket API and subscribe to the match channel to retrieve data for the `BTC-USD`, `ETH-USD` and `ETH-BTC` trading pairs.

The calculated VWAP of the respective trading pairs will be streamed to `STDOUT` as data from websocket connection becomes available.

### Makefile

A [Makefile](Makefile) is provided to bootstrap common CLI commands.

**Run project:**

```sh
make run
# go run .
```

**Test project:**

```sh
make test
# go test -v ./...
```
