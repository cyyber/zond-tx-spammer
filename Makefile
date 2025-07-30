# tx-spammer
VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS += -X 'github.com/theQRL/qrl-tx-spammer/utils.BuildVersion="$(VERSION)"'
GOLDFLAGS += -X 'github.com/theQRL/qrl-tx-spammer/utils.BuildTime="$(BUILDTIME)"'
GOLDFLAGS += -X 'github.com/theQRL/qrl-tx-spammer/utils.BuildRelease="$(RELEASE)"'

.PHONY: all test clean

all: build

test:
	go test ./...

build:
	@echo version: $(VERSION)
	go build -v -o bin/ -ldflags="-s -w $(GOLDFLAGS)" ./cmd/tx-spammer

clean:
	rm -f bin/*
