start:
	go run main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out