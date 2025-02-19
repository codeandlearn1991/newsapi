ARG GO_VERSION=1.24.0

# Build stage
FROM golang:${GO_VERSION}-bullseye AS builder
ARG APP

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/${APP}

# Runtime
FROM ubuntu:25.04

WORKDIR /app

COPY --from=builder /app .

EXPOSE 8080
CMD ["./app"]