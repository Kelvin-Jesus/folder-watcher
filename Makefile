build:
	@go build -o bin/goWatchSome

run: format build
	@./bin/goWatchSome

test:
	@go test -v ./...

format:
	@go fmt ./...