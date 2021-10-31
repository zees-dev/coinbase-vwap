# Coinbase Realtime Volume-Weighted-Average-Price Engine

A realtime VWAP (volume-weighted average price) calculation engine using the trading pair data from the [coinbase websocket](https://docs.pro.coinbase.com/#the-matches-channel).

This calculates the [VWAP](https://en.wikipedia.org/wiki/Volume-weighted_average_price) of crypto coin trading pairs as the data from the [match channel](https://docs.cloud.coinbase.com/exchange/docs/channels#match) becomes available (using websockets).

An arbitrary no. of trading pairs could be supported (by adding more to the [list](./main.go#L19)); however in our case, the following 3 trading pairs are defined:

- `BTC-USD`
- `ETH-USD`
- `ETH-BTC`

Note: it would be trivial to specify these trading pairs (along with the sliding window size) as CLI arguments.

## Key specs

- Support for specifying arbitrary no. of trading pairs (default to the 3 above)
- The VWAP will only be applicable to the first **200** data points per trading pair (configurable)
- The output will be streamed to a configurable `io.Writer` (`os.Stdout` in our case)
- Generic websocket client implementation - to which a [PubSuber](pkg/websocket/websocket.go#L4) can read from and write to (to communicate with the server)
  - Look at [tests](pkg/websocket/websocket_client_test.go) for more details
  - The design is extensible since different `PubSuber` implementations could be used to initate different types of subscriptions (i.e. to different coinbase channels, etc.) by communicating with the connected websocket
- Use of the [decimal](https://github.com/shopspring/decimal) for high-precision calculations (a must in finance)
- High test coverage

## Assumptions

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
