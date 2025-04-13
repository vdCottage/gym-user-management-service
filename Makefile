.PHONY: dev swagger build run clean

dev:
	~/go/bin/air

swagger:
	swag init -g cmd/server/main.go

build:
	go build -o ./tmp/main ./cmd/server

run:
	./tmp/main

clean:
	rm -rf ./tmp/

install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

.DEFAULT_GOAL := dev 