VERSION?=dev
COMMIT=$(shell git rev-parse HEAD | cut -c -8)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT}"
MODFLAGS=-mod=vendor
TESTFLAGS=-cover

PLATFORMS:=darwin linux windows

all: dev

dev:
	@echo building dist/zg...
	@go build ${MODFLAGS} ${LDFLAGS} -o dist/zg ./cmd/zg

dist: $(PLATFORMS)

$(PLATFORMS):
	@echo building dist/zg-$@-amd64...
	@GOOS=$@ GOARCH=amd64 go build ${MODFLAGS} ${LDFLAGS} -o dist/zg-$@-amd64 ./cmd/zg

clean:
	@rm -fr dist/

test:
	@go test ${MODFLAGS} ${TESTFLAGS} ./...

.PHONY: all dev dist $(PLATFORMS) clean test
