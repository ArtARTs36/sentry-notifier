lint:
	golangci-lint run --fix

test:
	go test ./...
