GO_BIN?=$(shell pwd)/.bin
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)

# Format the code
fmt::
	golangci-lint run --fix -v ./...

# Run the generate command
generate::
	go generate ./...

# Run the server
run::
	go run ./cmd/api-server/main.go

# Run test
test::
	go test ./...

# Run tidy
tidy::
	go mod tidy -v

tools::
	mkdir -p ${GO_BIN}
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_BIN} v1.61.0
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % sh -c 'GOBIN=${GO_BIN} go install %'
