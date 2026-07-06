package sdk

import (
	"context"
	"net"

	"github.com/Hell077/gale/internal/message"
	"github.com/Hell077/gale/internal/session"
)

const (
	SendSuccess = message.Success
	SendError   = message.Error
)

type Message = message.Message
type SendLog = message.SendLog

type Client struct {
	pool *session.ConnPool
	conn *session.Conn
}

var _ ClientAPI = (*Client)(nil)

func Connect(addr string) (*Client, error) {
	return ConnectContext(context.Background(), addr)
}

func ConnectContext(ctx context.Context, addr string) (*Client, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	if err := message.WriteHandshake(conn); err != nil {
		conn.Close()
		return nil, err
	}

	pool := session.NewPool()
	return &Client{
		pool: pool,
		conn: pool.AddTCPConnection(conn),
	}, nil
}

func NewMessage(topic uint64, data []byte) *Message {
	return message.NewMessage(topic, data)
}

func (c *Client) ID() uint64 {
	return c.conn.ID()
}

func (c *Client) Send(topic uint64, data []byte) (SendLog, error) {
	return c.SendMessage(NewMessage(topic, data))
}

func (c *Client) SendMessage(msg *Message) (SendLog, error) {
	return c.conn.SendMessage(msg)
}

func (c *Client) Receive() (*Message, error) {
	return c.conn.ReceiveMessage()
}

func (c *Client) ReceiveTopic(topic uint64) (*Message, error) {
	for {
		msg, err := c.Receive()
		if err != nil {
			return nil, err
		}
		if msg.Topic == topic {
			return msg, nil
		}
	}
}

func (c *Client) Close() error {
	return c.pool.RemoveConnection(c.conn.ID())
}
