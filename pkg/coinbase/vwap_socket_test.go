package coinbase

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zees-dev/coinbase-vwap/pkg/websocket"
)

func Test_SatisfiesPubSuberInterface(t *testing.T) {
	is := assert.New(t)
	is.Implements((*websocket.PubSuber)(nil), NewVWAPSocket(os.Stdout))
}

// testWriter is a test implementation of the io.Writer interface
type testWriter struct {
	buff []string
}

func (t *testWriter) Write(p []byte) (n int, err error) {
	t.buff = append(t.buff, string(p))
	return len(t.buff), nil
}

func Test_vwapSocket(t *testing.T) {
	is := assert.New(t)

	t.Run("initial subscription", func(t *testing.T) {
		vwapSocket := NewVWAPSocket(&testWriter{})
		go vwapSocket.ComputeVWAP(context.Background(), 2, "BTC-USD")

		initialSub := <-vwapSocket.Send()
		msg, err := toMessage(initialSub)
		is.NoError(err)
		is.Equal(string(msg.Type), "subscribe")
	})

	t.Run("writer writes vwap computation with single pair", func(t *testing.T) {
		writer := &testWriter{}
		vwapSocket := NewVWAPSocket(writer)
		go vwapSocket.ComputeVWAP(context.Background(), 2, "BTC-USD")

		// example websocket message
		vwapSocket.Recv() <- []byte(`{"type":"match","size":"0.01168318","price":"4275.38","product_id":"BTC-USD"}`)

		// give the socket time to recv the message
		<-time.After(time.Millisecond * 100)

		is.Equal("vwap BTC-USD: 4275.38\n", writer.buff[0])
	})

	t.Run("writer writes vwap computation with multiple pairs", func(t *testing.T) {
		writer := &testWriter{}
		vwapSocket := NewVWAPSocket(writer)
		go vwapSocket.ComputeVWAP(context.Background(), 2, "BTC-USD", "ETH-USD")

		// example websocket messages
		vwapSocket.Recv() <- []byte(`{"type":"match","size":"0.01168318","price":"4275.38","product_id":"BTC-USD"}`)
		vwapSocket.Recv() <- []byte(`{"type":"match","size":"0.01168318","price":"4000.38","product_id":"ETH-USD"}`)

		// give the socket time to recv the messages
		<-time.After(time.Millisecond * 100)

		is.Equal("vwap BTC-USD: 4275.38\n", writer.buff[0])
		is.Equal("vwap ETH-USD: 4000.38\n", writer.buff[1])
	})

}
