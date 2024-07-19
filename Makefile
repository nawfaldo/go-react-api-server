build:
	@go build -o bin/test cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/test