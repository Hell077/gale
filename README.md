<p align="center">
  <img src="./assets/gale.png" width="350" alt="Gale logo">
</p>

# Gale Mini Guide

Gale is a small TCP message transport for Go processes, servers, and clusters. It sends binary framed messages with a `topic`, `offset`, and `data` payload.

## Install

Use the SDK from another Go project:

```bash
go get github.com/Hell077/gale
```

Import it:

```go
import "github.com/Hell077/gale/sdk"
```

## Run Gale as a Binary

Run the TCP server directly:

```bash
go run ./cmd -addr 0.0.0.0:7827
```

Or configure it with environment variables:

```bash
GALE_HOST=0.0.0.0 GALE_PORT=7827 go run ./cmd
```

`GALE_ADDR` can be used when you want to pass the full listen address:

```bash
GALE_ADDR=127.0.0.1:7828 go run ./cmd
```

Build a local binary:

```bash
go build -o bin/gale ./cmd
./bin/gale -addr 0.0.0.0:7827
```

The server listens until it receives `SIGINT` or `SIGTERM`. On shutdown it closes the TCP listener and active connections.

## Docker

The Dockerfile is located at the repository root:

```text
./Dockerfile
```

Build the image:

```bash
docker build -t gale .
```

Run it:

```bash
docker run --rm -p 7827:7827 gale
```

Use a custom port:

```bash
docker run --rm -e GALE_PORT=7828 -p 7828:7828 gale
```

Pull the published image from Docker Hub:

```bash
docker pull alexeyovchinnikovhell077/gale:latest
docker run --rm -p 7827:7827 alexeyovchinnikovhell077/gale:latest
```

The GitHub Actions workflow in `.github/workflows/dockerhub.yml` builds this
Dockerfile when `Dockerfile`, `.dockerignore`, `cmd/**`, `internal/**`,
`sdk/**`, or Go module files change. Pull requests only build the image for
validation. Pushes and manual workflow runs publish the image to Docker Hub.

Configure these repository secrets before publishing:

- `DOCKERHUB_USERNAME`: `alexeyovchinnikovhell077`.
- `DOCKERHUB_TOKEN`: Docker Hub access token with push access.

Optional repository variable:

- `DOCKERHUB_REPOSITORY`: Docker Hub repository name. Defaults to `gale`.

Configuration order:

1. `-addr`
2. `GALE_ADDR`
3. `GALE_HOST` and `GALE_PORT`
4. default `0.0.0.0:7827`

## Docker Compose Integration

Add Gale to an existing `docker-compose.yml`:

```yaml
services:
  gale:
    image: gale:local
    build:
      context: ./gale
    environment:
      GALE_HOST: 0.0.0.0
      GALE_PORT: 7827
    ports:
      - "7827:7827"
```

If you publish Gale as an image, use it without `build`:

```yaml
services:
  gale:
    image: alexeyovchinnikovhell077/gale:latest
    environment:
      GALE_HOST: 0.0.0.0
      GALE_PORT: 7827
    ports:
      - "7827:7827"
```

Other services in the same compose network can connect to `gale:7827`.

## SDK Server

Start a server from Go code:

```go
package main

import (
	"log"

	"github.com/Hell077/gale/sdk"
)

func main() {
	server, err := sdk.Listen("127.0.0.1:7827")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Printf("listening on %s", server.Addr())

	select {}
}
```

## SDK Client

Connect and send a message:

```go
package main

import (
	"log"

	"github.com/Hell077/gale/sdk"
)

func main() {
	client, err := sdk.Connect("127.0.0.1:7827")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if _, err := client.Send(1, []byte("hello")); err != nil {
		log.Fatal(err)
	}
}
```

## Receiving Messages

Receive any message:

```go
msg, err := client.Receive()
if err != nil {
	log.Fatal(err)
}
log.Printf("topic=%d data=%s", msg.Topic, msg.Data)
```

Receive only one topic:

```go
msg, err := client.ReceiveTopic(1)
if err != nil {
	log.Fatal(err)
}
log.Printf("topic=%d data=%s", msg.Topic, msg.Data)
```

`ReceiveTopic(1)` skips messages from other topics and returns the next message from topic `1`.

## Server Connections

The server can inspect active connections and send messages back:

```go
for _, conn := range server.Connections() {
	if _, err := conn.Send(1, []byte("server message")); err != nil {
		log.Printf("send to %d failed: %v", conn.ID(), err)
	}
}
```

Broadcast to all active connections:

```go
logs := server.Broadcast(1, []byte("broadcast message"))
log.Printf("sent to %d connections", len(logs))
```

## SDK Interfaces

The SDK exposes small interfaces for testing and integration:

```go
type Sender interface {
	Send(topic uint64, data []byte) (sdk.SendLog, error)
	SendMessage(msg *sdk.Message) (sdk.SendLog, error)
}
```

Common interfaces:

- `sdk.ClientAPI`
- `sdk.ServerAPI`
- `sdk.ConnectionAPI`
- `sdk.Sender`
- `sdk.Receiver`
- `sdk.Broadcaster`
