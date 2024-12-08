# News API Server - A Sample Go REST API

This project demonstrates the development of a RESTful API service using the Go programming language. It aims to provide a solid foundation for building your own APIs, adhering to common patterns and conventions.

## Features

- RESTful API design following best practices (GET, POST, PUT, DELETE for CRUD operations)
- Modular code structure with clear separation of concerns (handlers, models, database access, etc.)
- Built-in unit tests for core functionalities

## Installation

Prerequisites:

Go version 1.x or later (check with go version)
A code editor or IDE with Go support

Install dependencies:

`go mod tidy`

## Running the API

Run the API:

```sh
go run ./cmd/api-server
```

This will start the API server on port 8080 by default (adjust the port if needed).

## API Endpoints

POST /news - Create a new news resource
GET /news - Retrieve a list of all news
GET /news/:id - Get details of a specific news by ID
PUT /news/:id - Update an existing news
DELETE /news/:id - Delete a news

## Testing

Unit tests are crucial for ensuring code quality. You can run your tests with:

```sh
go test ./...
```

## License

This project is licensed under the MIT License (see LICENSE file).
