# Sch Code Review Repository

This project is built with Go and uses `swaggo` to generate Swagger documentation for API endpoints. The following instructions will guide you through setting up, linting, testing, and building the project.

## Getting Started

### Prerequisites

Ensure you have the following dependencies installed:

- Go 1.22+
- golangci-lint
- swaggo

### Installation

To set up the project, run the following commands to install necessary packages:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### Generate Swagger Documentation

Generate Swagger documentation based on the API code:

```bash
swag init -g cmd/coupon_service/main.go
```

### Build the Project

To lint, test, and build the project, you can use the following commands:

```bash
make
```

Or you can execute the individual targets as described below.

## Makefile Commands

### Lint

To lint the project, run:

```bash
make lint
```

To automatically fix linting issues, run:

```bash
make lint-fix
```

### Tests

Run the tests with:

```bash
make test
```

### Build

Build the project with:

```bash
make build
```

## Testing and Code Coverage

To run all tests with verbose output, use:

```bash
go test -v ./...
```

For code coverage details, run:

```bash
go test -v -cover ./...
```

Current code coverage is **83.3%** of statements.

## Docker

To build a Docker image, run:

```bash
docker build -t schwarz-it-code-review .
```

To run the Docker container, use:

```bash
docker run -p 8080:8080 schwarz-it-code-review
```

The API will be available at `http://localhost:8080/swagger/index.html`.

## Docker Compose

To run the application with Docker Compose, use:

Up:
```bash
docker-compose up
```

Down:
```bash
docker-compose down
```

Build and run the application with Docker Compose:

```bash
docker-compose up --build
```
