package main

import "testing"

func TestResolveAddr(t *testing.T) {
	t.Setenv("GALE_ADDR", "")
	t.Setenv("GALE_HOST", "")
	t.Setenv("GALE_PORT", "")

	if got := resolveAddr("127.0.0.1:7000"); got != "127.0.0.1:7000" {
		t.Fatalf("expected flag addr, got %q", got)
	}

	t.Setenv("GALE_ADDR", "127.0.0.1:7001")
	if got := resolveAddr(""); got != "127.0.0.1:7001" {
		t.Fatalf("expected GALE_ADDR, got %q", got)
	}

	t.Setenv("GALE_ADDR", "")
	t.Setenv("GALE_HOST", "127.0.0.1")
	t.Setenv("GALE_PORT", "7002")
	if got := resolveAddr(""); got != "127.0.0.1:7002" {
		t.Fatalf("expected host and port addr, got %q", got)
	}

	t.Setenv("GALE_HOST", "")
	t.Setenv("GALE_PORT", "")
	if got := resolveAddr(""); got != "0.0.0.0:9000" {
		t.Fatalf("expected default addr, got %q", got)
	}
}
