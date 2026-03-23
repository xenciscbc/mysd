BINARY=mysd
VERSION?=dev

.PHONY: build test lint clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./... -v -count=1

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
