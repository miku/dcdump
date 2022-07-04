SHELL := /bin/bash
TARGETS := dcdump
# VERSION := $(shell git rev-parse --short HEAD)
VERSION := "v0.1.3"
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOLDFLAGS += -w -s
GOFLAGS = -ldflags "$(GOLDFLAGS)"

.PHONY: all
all: $(TARGETS)

%: cmd/%/main.go
	go build -o $@ -ldflags "$(GOLDFLAGS)" $<

.PHONY: clean
clean:
	rm -f $(TARGETS)

.PHONY: update-all-deps
update-all-deps:
	go get -u -v ./... && go mod tidy

