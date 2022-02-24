.PHONY: build
build: 
	go build -trimpath -ldflags=-w -o store-location

.PHONY: test
test:

.PHONY: sanity-check
sanity-check:
	golangci-lint run

.PHONY: test-unit
test-unit:
	go test -v ./...
