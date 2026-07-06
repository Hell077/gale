package sdk

import (
	"context"

	"github.com/Hell077/gale/internal/dispatcher"
	"github.com/Hell077/gale/internal/session"
)

type Server struct {
	dispatcher *dispatcher.Dispatcher
	cancel     context.CancelFunc
}

type Connection struct {
	conn *session.Conn
}

var _ ServerAPI = (*Server)(nil)
var _ ConnectionAPI = (*Connection)(nil)

func Listen(addr string) (*Server, error) {
	return ListenContext(context.Background(), addr)
}

func ListenContext(ctx context.Context, addr string) (*Server, error) {
	ctx, cancel := context.WithCancel(ctx)
	d := dispatcher.NewDispatcher()
	if err := d.ListenTCP(ctx, addr); err != nil {
		cancel()
		return nil, err
	}
	return &Server{
		dispatcher: d,
		cancel:     cancel,
	}, nil
}

func (s *Server) Addr() string {
	return s.dispatcher.TCPAddr()
}

func (s *Server) Connections() []Connection {
	conns := s.dispatcher.Pool.Snapshot()
	result := make([]Connection, 0, len(conns))
	for _, conn := range conns {
		result = append(result, Connection{conn: conn})
	}
	return result
}

func (s *Server) Broadcast(topic uint64, data []byte) []SendLog {
	return s.BroadcastMessage(NewMessage(topic, data))
}

func (s *Server) BroadcastMessage(msg *Message) []SendLog {
	return s.dispatcher.Pool.Broadcast(msg)
}

func (s *Server) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return s.dispatcher.Close()
}

func (c Connection) ID() uint64 {
	return c.conn.ID()
}

func (c Connection) Send(topic uint64, data []byte) (SendLog, error) {
	return c.SendMessage(NewMessage(topic, data))
}

func (c Connection) SendMessage(msg *Message) (SendLog, error) {
	return c.conn.SendMessage(msg)
}

func (c Connection) Receive() (*Message, error) {
	return c.conn.ReceiveMessage()
}

func (c Connection) ReceiveTopic(topic uint64) (*Message, error) {
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

func (c Connection) Close() error {
	return c.conn.Close()
}
