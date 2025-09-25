GO ?= go
GOFMT ?= gofmt "-s"
PLATFORM ?=darwin
GOARCH ?=arm64
# GOARM ?= 7
GOFILES := $(shell find . -name "*.go" -type f)

all: install

install:
	env GOOS=$(PLATFORM) GOARCH=$(GOARCH) GOARM=$(GOARM) $(GO) build -o ./serpcli ./
	mv ./serpcli /usr/local/bin/serpcli

fmt:
	$(GOFMT) -w $(GOFILES)

test:
	go test -cover ./...