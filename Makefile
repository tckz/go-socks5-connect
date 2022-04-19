.PHONY: all lint clean test

export GO111MODULE=on

ifeq ($(GO_CMD),)
GO_CMD:=go
endif

VERSION := $(shell git describe --always)
GO_BUILD := CGO_ENABLED=0 $(GO_CMD) build -ldflags "-X main.version=$(VERSION)"

DIST_SOCKS5_CONNECT=dist/socks5-connect

TARGETS = \
	$(DIST_SOCKS5_CONNECT)

all: $(TARGETS)
	@echo "$@ done." 1>&2

clean:
	/bin/rm -f $(TARGETS)
	@echo "$@ done." 1>&2

$(DIST_SOCKS5_CONNECT): cmd/socks5-connect/*
	$(GO_BUILD) -o $@ ./cmd/socks5-connect/

