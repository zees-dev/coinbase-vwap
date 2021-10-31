package websocket

// PubSuber represents an interface which allows communication (publish/subscribe) to/from a websocket.
type PubSuber interface {
	Send() chan []byte
	Recv() chan []byte
}
