package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Hell077/gale/internal/dispatcher"
)

func main() {
	addr := flag.String("addr", "", "tcp address to listen on")
	flag.Parse()

	listenAddr := resolveAddr(*addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	d := dispatcher.NewDispatcher()
	if err := d.ListenTCP(ctx, listenAddr); err != nil {
		log.Fatalf("listen tcp: %v", err)
	}

	log.Printf("gale listening on %s", d.TCPAddr())

	<-ctx.Done()
	log.Printf("shutdown signal received")

	if err := d.Close(); err != nil {
		log.Printf("shutdown error: %v", err)
		os.Exit(1)
	}
	log.Printf("gale stopped")
}

func resolveAddr(flagAddr string) string {
	if flagAddr != "" {
		return flagAddr
	}
	if envAddr := strings.TrimSpace(os.Getenv("GALE_ADDR")); envAddr != "" {
		return envAddr
	}

	host := strings.TrimSpace(os.Getenv("GALE_HOST"))
	if host == "" {
		host = "0.0.0.0"
	}

	port := strings.TrimSpace(os.Getenv("GALE_PORT"))
	if port == "" {
		port = "7827"
	}

	return host + ":" + port
}
