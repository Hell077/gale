APP ?= gale
HOST ?= 0.0.0.0
PORT ?= 9000
ADDR ?= $(HOST):$(PORT)
IMAGE ?= gale:local
BIN_DIR ?= bin
GOCACHE ?= /private/tmp/gale-go-cache
GOMODCACHE ?= /private/tmp/gale-go-mod

.PHONY: help test build run clean docker-build docker-run docker-shell

help:
	@echo "Targets:"
	@echo "  test          Run Go tests"
	@echo "  build         Build $(APP) into $(BIN_DIR)/$(APP)"
	@echo "  run           Run $(APP) locally with ADDR=$(ADDR)"
	@echo "  clean         Remove local build output"
	@echo "  docker-build  Build Docker image IMAGE=$(IMAGE)"
	@echo "  docker-run    Run Docker image on ADDR=$(ADDR)"
	@echo "  docker-shell  Open a shell in the Docker image"

test:
	GOCACHE=$(GOCACHE) go test ./...

build:
	mkdir -p $(BIN_DIR)
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build -o $(BIN_DIR)/$(APP) ./cmd

run:
	GOCACHE=$(GOCACHE) go run ./cmd -addr $(ADDR)

clean:
	rm -rf $(BIN_DIR)

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm -p $(PORT):$(PORT) -e GALE_HOST=0.0.0.0 -e GALE_PORT=$(PORT) $(IMAGE)

docker-shell:
	docker run --rm -it --entrypoint sh $(IMAGE)
