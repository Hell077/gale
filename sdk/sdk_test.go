package sdk

import (
	"testing"
	"time"
)

func TestSDKClientServerRoundTrip(t *testing.T) {
	server, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer server.Close()

	client, err := Connect(server.Addr())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer client.Close()

	var conn Connection
	deadline := time.After(time.Second)
	for conn.conn == nil {
		conns := server.Connections()
		if len(conns) > 0 {
			conn = conns[0]
			break
		}

		select {
		case <-deadline:
			t.Fatal("server did not accept client")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	if _, err := client.Send(7, []byte("hello from sdk")); err != nil {
		t.Fatalf("client send: %v", err)
	}

	got, err := conn.Receive()
	if err != nil {
		t.Fatalf("server receive: %v", err)
	}
	if got.Topic != 7 || string(got.Data) != "hello from sdk" {
		t.Fatalf("unexpected client message: topic=%d data=%q", got.Topic, got.Data)
	}

	if _, err := conn.Send(8, []byte("hello back")); err != nil {
		t.Fatalf("server send: %v", err)
	}

	got, err = client.Receive()
	if err != nil {
		t.Fatalf("client receive: %v", err)
	}
	if got.Topic != 8 || string(got.Data) != "hello back" {
		t.Fatalf("unexpected server message: topic=%d data=%q", got.Topic, got.Data)
	}
}

func TestSDKReceiveTopicSkipsOtherTopics(t *testing.T) {
	server, err := Listen("127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer server.Close()

	client, err := Connect(server.Addr())
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer client.Close()

	conn := waitSDKConnection(t, server)

	if _, err := conn.Send(2, []byte("skip 2")); err != nil {
		t.Fatalf("send topic 2: %v", err)
	}
	if _, err := conn.Send(3, []byte("skip 3")); err != nil {
		t.Fatalf("send topic 3: %v", err)
	}
	if _, err := conn.Send(1, []byte("take 1")); err != nil {
		t.Fatalf("send topic 1: %v", err)
	}

	got, err := client.ReceiveTopic(1)
	if err != nil {
		t.Fatalf("receive topic 1: %v", err)
	}
	if got.Topic != 1 || string(got.Data) != "take 1" {
		t.Fatalf("unexpected topic message: topic=%d data=%q", got.Topic, got.Data)
	}
}

func waitSDKConnection(t *testing.T, server ServerAPI) Connection {
	t.Helper()

	deadline := time.After(time.Second)
	for {
		conns := server.Connections()
		if len(conns) > 0 {
			return conns[0]
		}

		select {
		case <-deadline:
			t.Fatal("server did not accept client")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
