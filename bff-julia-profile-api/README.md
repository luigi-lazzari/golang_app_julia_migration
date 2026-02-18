# BFF Julia Profile API - Go Implementation

Backend For Frontend (BFF) API for Julia User Profile Management implemented in Go.

## Features

- RESTful API with Gin framework
- User preferences management
- User profile management
- Azure App Configuration integration
- Azure Cosmos DB for data persistence
- OpenAPI/Swagger documentation
- Prometheus metrics
- Health check endpoints
- Structured logging with Zap

## Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Azure App Configuration Emulator
- Azure Cosmos DB Emulator

## Quick Start

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go

# Or use make
make run
```

### Docker

```bash
# Build and run with docker-compose
docker-compose up --build -d

# View logs
docker-compose logs -f
```

## API Documentation

- Swagger UI: http://localhost:8090/swagger/index.html
- OpenAPI JSON: http://localhost:8090/swagger/doc.json

## Health & Metrics

- Health: http://localhost:8090/health
- Metrics: http://localhost:8090/metrics

## API Endpoints

### User Preferences
- `GET /api/v1/user/preferences` - Get user preferences
- `PUT /api/v1/user/preferences` - Update user preferences

### User Profile
- `GET /api/v1/user/profile` - Get user profile
- `PUT /api/v1/user/profile` - Update user profile

## Configuration

Configuration is loaded from:
1. Environment variables
2. Azure App Configuration (local emulator by default)
3. `config.yaml` file

### Environment Variables

```bash
SERVER_PORT=8090
AZURE_APPCONFIG_ENDPOINT=http://localhost:8484
COSMOS_DB_ENDPOINT=https://localhost:8182
COSMOS_DB_KEY=<emulator-key>
COSMOS_DB_DATABASE=bff_julia_db
LOG_LEVEL=info
```

## Project Structure

```
bff-julia-profile-api/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── handler/              # HTTP handlers (controllers)
│   ├── service/              # Business logic
│   ├── repository/           # Data access layer
│   ├── model/                # Domain models
│   └── middleware/           # HTTP middlewares
├── pkg/
│   ├── azure/                # Azure SDK utilities
│   └── logger/               # Logging utilities
├── api/
│   └── openapi.yaml          # OpenAPI specification
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Building

```bash
# Build binary
make build

# Run tests
make test

# Generate OpenAPI docs
make swagger
```

## License

Comune di Roma
# golang_app_julia_migration
