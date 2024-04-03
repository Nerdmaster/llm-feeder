.PHONY: build

build:
	go build -o ./bin/feeder .

lint:
	revive ./...
	go vet ./...
