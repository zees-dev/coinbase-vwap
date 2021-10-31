package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type wsClient struct {
	conn *websocket.Conn
	wsIO PubSuber
	lock sync.Mutex
}

// NewWSClient initialises client and connects to websocket server at specified URL.
func NewWSClient(url string, wsIO PubSuber) (*wsClient, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "dialing")
	}
	log.Printf("successfully connected to %s...", url)
	return &wsClient{conn: ws, wsIO: wsIO}, nil
}

// Start initialised communication between sender and reciever channels.
// The context can be used to optionally cancel/close the connection.
// note: this is a blocking operation
func (c *wsClient) Start(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// non-blocking message recieval from server
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read err:", err)
				return
			}
			// log.Printf("recv-raw: %s\n", msg)
			c.wsIO.Recv() <- msg
		}
	}()

	// handle all scenarios
	for {
		select {
		case <-done:
			log.Println("server terminated the connection")
			return nil
		case req := <-c.wsIO.Send(): // process messages from sender channel
			// log.Printf("send: %s\n", req)
			err := c.conn.WriteMessage(websocket.TextMessage, req)
			if err != nil {
				return errors.Wrap(err, "error sending message")
				// log.Println("error sending message:", err)
				// return
			}
		case <-ctx.Done(): // process OS interrupts (or context cancellation)
			log.Println("interrupt (context cancelled)")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return errors.Wrap(err, "error closing websocket connection")
				// log.Println("error closing websocket connection:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
