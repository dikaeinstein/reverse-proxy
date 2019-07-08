BINARY_NAME=reverse-proxy

## Fetch dependencies
install:
	go get -t -v ./...

## Run tests
test:
	APP_ENV=test go test -race -cover -v -coverprofile=c.out ./...

coverage:
	go tool cover -html=c.out -o coverage.html

## Build binary
build:
	go build

## Execute binary
run:build
	./$(BINARY_NAME)

.PHONY: clean
## Remove binary
clean:
	rm -f $(BINARY_NAME)
