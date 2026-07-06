package session

import (
	"context"
	"errors"
	"net"

	"github.com/Hell077/gale/internal/message"
)

type TCPServer struct {
	listener net.Listener
	pool     *ConnPool
}

func ListenTCP(ctx context.Context, addr string, pool *ConnPool) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		pool = NewPool()
	}

	server := &TCPServer{
		listener: listener,
		pool:     pool,
	}

	go server.accept(ctx)
	return server, nil
}

func (s *TCPServer) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *TCPServer) Pool() *ConnPool {
	return s.pool
}

func (s *TCPServer) Close() error {
	return s.listener.Close()
}

func (s *TCPServer) accept(ctx context.Context) {
	go func() {
		<-ctx.Done()
		_ = s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				return
			}
			return
		}

		go func(conn net.Conn) {
			if err := message.ReadHandshake(conn); err != nil {
				_ = conn.Close()
				return
			}
			s.pool.AddTCPConnection(conn)
		}(conn)
	}
}
