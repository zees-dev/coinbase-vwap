package websocket

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/zees-dev/coinbase-vwap/pkg/vwap"
)

type vwapSocket struct {
	writer io.Writer   // writer to write output to
	req    chan []byte // request channel to send all requests from caller
	res    chan []byte // response channel to receive all responses from websocket
	lock   sync.Mutex
}

func NewVWAPSocket(writer io.Writer) *vwapSocket {
	// channel for requests - we will use this to send messages to the websocket
	// channel for responses - we will use this to receive messages from the websocket
	requestCh, responseCh := make(chan []byte), make(chan []byte)
	return &vwapSocket{writer: writer, req: requestCh, res: responseCh}
}

// Send requests to websocket on this channel
func (v *vwapSocket) Send() chan []byte {
	return v.req
}

// Recv response from websocket on this channel
func (v *vwapSocket) Recv() chan []byte {
	return v.res
}

// Connect will interop with the websocket client via utilising the same IO channels.
// vwapSocket implements the PubSuber interface.
// note: this is a blocking operation
func (v *vwapSocket) ComputeVWAP(ctx context.Context, slidingWindowSize int, pairs ...string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	// create vwap counters for each pair
	vwapMap := make(map[string]vwap.Vwapper)
	for _, pair := range pairs {
		vwapMap[pair] = vwap.NewVWAPCounter(slidingWindowSize)
	}

	// initial subscription - send once consumer is ready (websocket connection is established)
	go func() {
		initialSub := subscription{
			Type:     subscribe,
			Channels: []channel{{Name: matches, ProductIDs: pairs}},
		}

		v.Send() <- initialSub.Byte()
		log.Printf("send: %+v", initialSub)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("context cancelled")
			return nil
		case res := <-v.Recv():
			msg, err := toMessage(res)
			if err != nil {
				log.Printf("error: %s\n", err)
			}
			log.Printf("recv: %s\n", res)

			vwapMapCounter, ok := vwapMap[msg.ProductID]
			if !ok {
				log.Printf("unexpected product id: %s\n", msg.ProductID)
				continue
			}

			vwapMapCounter.Update(msg.Price, msg.Size)
			outText := fmt.Sprintf("vwap %s: %s\n", msg.ProductID, vwapMapCounter.VWAP())
			v.writer.Write([]byte(outText))
		}
	}
}
