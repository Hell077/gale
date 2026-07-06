package session

import (
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/Hell077/gale/internal/message"
)

const size uint16 = 32 * 1024

type ConnPool struct {
	//connections counter
	counter atomic.Uint64

	//current connections
	mu          sync.RWMutex
	Connections map[uint64]*Conn

	messageStore map[int64]*message.Message
}

type Conn struct {
	internalID uint64

	//message queue
	queue Queue

	rw io.ReadWriteCloser
}

func NewPool() *ConnPool {
	return &ConnPool{
		Connections: make(map[uint64]*Conn),
	}
}

func (pool *ConnPool) CreateConnection() *Conn {
	pr, pw := io.Pipe()
	return pool.addConn(&pipeReadWriteCloser{
		reader: pr,
		writer: pw,
	})
}

func (pool *ConnPool) AddTCPConnection(conn net.Conn) *Conn {
	return pool.addConn(conn)
}

func (pool *ConnPool) DialTCP(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	if err := message.WriteHandshake(conn); err != nil {
		conn.Close()
		return nil, err
	}
	return pool.addConn(conn), nil
}

func (pool *ConnPool) addConn(rw io.ReadWriteCloser) *Conn {
	conn := &Conn{
		internalID: pool.counter.Add(1),
		rw:         rw,
	}
	pool.mu.Lock()
	pool.Connections[conn.internalID] = conn
	pool.mu.Unlock()

	return conn
}

func (pool *ConnPool) RemoveConnection(id uint64) error {
	pool.mu.Lock()
	conn, ok := pool.Connections[id]
	if ok {
		delete(pool.Connections, id)
	}
	pool.mu.Unlock()
	if !ok {
		return nil
	}
	return conn.Close()
}

func (pool *ConnPool) Close() error {
	pool.mu.Lock()
	conns := make([]*Conn, 0, len(pool.Connections))
	for id, conn := range pool.Connections {
		conns = append(conns, conn)
		delete(pool.Connections, id)
	}
	pool.mu.Unlock()

	var closeErr error
	for _, conn := range conns {
		if err := conn.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	return closeErr
}

func (pool *ConnPool) Broadcast(msg *message.Message) []message.SendLog {
	pool.mu.RLock()
	conns := make([]*Conn, 0, len(pool.Connections))
	for _, conn := range pool.Connections {
		conns = append(conns, conn)
	}
	pool.mu.RUnlock()

	logs := make([]message.SendLog, 0, len(conns))
	for _, conn := range conns {
		log, _ := conn.SendMessage(msg)
		logs = append(logs, log)
	}
	return logs
}

func (pool *ConnPool) Snapshot() []*Conn {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	conns := make([]*Conn, 0, len(pool.Connections))
	for _, conn := range pool.Connections {
		conns = append(conns, conn)
	}
	return conns
}

func (s *Conn) ID() uint64 {
	return s.internalID
}

func (s *Conn) Close() error {
	return s.rw.Close()
}

type pipeReadWriteCloser struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func (p *pipeReadWriteCloser) Read(b []byte) (int, error) {
	return p.reader.Read(b)
}

func (p *pipeReadWriteCloser) Write(b []byte) (int, error) {
	return p.writer.Write(b)
}

func (p *pipeReadWriteCloser) Close() error {
	readErr := p.reader.Close()
	writeErr := p.writer.Close()
	if readErr != nil {
		return readErr
	}
	return writeErr
}
