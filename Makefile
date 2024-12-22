# Format the code
fmt::
	go fmt ./...

# Run the server
run::
	go run ./cmd/api-server/main.go

# Run tidy
tidy::
	go mod tidy -v