package dispatcher

import (
	"context"

	"github.com/Hell077/gale/internal/session"
)

type Dispatcher struct {
	Pool      *session.ConnPool
	tcpServer *session.TCPServer
	alive     bool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		Pool: session.NewPool(),
	}
}

func (d *Dispatcher) ListenTCP(ctx context.Context, addr string) error {
	server, err := session.ListenTCP(ctx, addr, d.Pool)
	if err != nil {
		return err
	}
	d.tcpServer = server
	d.alive = true
	return nil
}

func (d *Dispatcher) TCPAddr() string {
	if d.tcpServer == nil {
		return ""
	}
	return d.tcpServer.Addr().String()
}

func (d *Dispatcher) Close() error {
	d.alive = false
	var closeErr error
	if d.tcpServer == nil {
		return d.Pool.Close()
	}
	if err := d.tcpServer.Close(); err != nil {
		closeErr = err
	}
	if err := d.Pool.Close(); err != nil && closeErr == nil {
		closeErr = err
	}
	return closeErr
}
