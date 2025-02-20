build:
	@go build -o bin/phuntgraphvisualize

run: build
	@./bin/phuntgraphvisualize

test:
	go test -v ./...