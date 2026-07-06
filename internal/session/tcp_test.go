package session

import (
	"context"
	"testing"
	"time"

	"github.com/Hell077/gale/internal/message"
)

func TestTCPMessageRoundTrip(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverPool := NewPool()
	server, err := ListenTCP(ctx, "127.0.0.1:0", serverPool)
	if err != nil {
		t.Fatalf("listen tcp: %v", err)
	}
	defer server.Close()

	clientPool := NewPool()
	clientConn, err := clientPool.DialTCP(server.Addr().String())
	if err != nil {
		t.Fatalf("dial tcp: %v", err)
	}
	defer clientConn.Close()

	var serverConn *Conn
	deadline := time.After(time.Second)
	for serverConn == nil {
		conns := serverPool.Snapshot()
		if len(conns) > 0 {
			serverConn = conns[0]
			break
		}

		select {
		case <-deadline:
			t.Fatal("server did not accept tcp connection")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	want := message.NewMessage(42, []byte("cluster payload"))
	if _, err := clientConn.SendMessage(want); err != nil {
		t.Fatalf("send message: %v", err)
	}

	got, err := serverConn.ReceiveMessage()
	if err != nil {
		t.Fatalf("receive message: %v", err)
	}
	if got.Topic != want.Topic || string(got.Data) != string(want.Data) {
		t.Fatalf("unexpected message: topic=%d data=%q", got.Topic, got.Data)
	}
}