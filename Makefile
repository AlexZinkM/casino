run:
	swag init
	golangci-lint run
	go mod tidy
	go test ./...
	go run main.go
