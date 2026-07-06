package sdk

import "io"

type Sender interface {
	Send(topic uint64, data []byte) (SendLog, error)
	SendMessage(msg *Message) (SendLog, error)
}

type Receiver interface {
	Receive() (*Message, error)
	ReceiveTopic(topic uint64) (*Message, error)
}

type Peer interface {
	io.Closer
	Sender
	Receiver

	ID() uint64
}

type Broadcaster interface {
	Broadcast(topic uint64, data []byte) []SendLog
	BroadcastMessage(msg *Message) []SendLog
}

type ClientAPI interface {
	Peer
}

type ConnectionAPI interface {
	Peer
}

type ServerAPI interface {
	io.Closer
	Broadcaster

	Addr() string
	Connections() []Connection
}
