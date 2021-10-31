all: clean tidy test build run

clean:
	rm -rf wvap

tidy:
	go mod tidy

test:
	go test -v ./...
.PHONY: test

build:
	go build -o vwap main.go

run:
	go run .
