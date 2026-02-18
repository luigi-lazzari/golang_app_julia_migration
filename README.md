# Julia App - Go Projects

Questo repository contiene le implementazioni Go dei Backend For Frontend (BFF) per l'applicazione Julia del Comune di Roma.

## Progetti

### bff-julia-mobile-api
Backend For Frontend per l'applicazione mobile Julia. Gestisce la configurazione dell'app basata su piattaforma e versione.

**Porta:** 8080  
**Endpoints principali:**
- `GET /api/v1/app-config` - Recupera la configurazione dell'app

### bff-julia-profile-api
Backend For Frontend per la gestione del profilo utente e delle preferenze.

**Porta:** 8090  
**Endpoints principali:**
- `GET /api/v1/user/profile` - Recupera il profilo utente
- `PUT /api/v1/user/profile` - Aggiorna il profilo utente
- `GET /api/v1/user/preferences` - Recupera le preferenze utente
- `PUT /api/v1/user/preferences` - Aggiorna le preferenze utente

## Tecnologie

- **Go 1.23+**
- **Gin** - Web framework
- **Azure SDK for Go** - Integrazione con Azure Cosmos DB e App Configuration
- **Prometheus** - Metriche e monitoring
- **Zap** - Structured logging
- **Docker** - Containerizzazione

## Quick Start

### Prerequisiti

```bash
# Installa Go 1.23 o superiore
go version

# Installa Docker e Docker Compose
docker --version
docker-compose --version
```

### Avvio di entrambi i progetti

#### Mobile API
```bash
cd bff-julia-mobile-api
make run
# Oppure con Docker
make docker-run
```

#### Profile API
```bash
cd bff-julia-profile-api
make run
# Oppure con Docker
make docker-run
```

## Struttura del Progetto

Entrambi i progetti seguono la stessa architettura:

```
project/
├── cmd/api/              # Entry point dell'applicazione
├── internal/
│   ├── config/          # Gestione configurazione
│   ├── handler/         # HTTP handlers (REST endpoints)
│   ├── service/         # Business logic
│   ├── repository/      # Data access layer
│   ├── model/           # Domain models
│   └── middleware/      # HTTP middlewares
├── pkg/
│   ├── azure/           # Azure SDK utilities
│   └── logger/          # Logging utilities
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Configurazione

Le applicazioni utilizzano variabili d'ambiente per la configurazione:

```bash
# Server
SERVER_PORT=8080  # o 8090 per profile-api

# Azure App Configuration
AZURE_APPCONFIG_ENDPOINT=http://localhost:8484
AZURE_APPCONFIG_CONNECTION_STRING=Endpoint=http://localhost:8484;Id=local;Secret=c2VjcmV0

# Azure Cosmos DB
COSMOS_DB_ENDPOINT=https://localhost:8182
COSMOS_DB_KEY=C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==
COSMOS_DB_DATABASE=bff_julia_db
COSMOS_EMULATOR_ENABLED=true

# Logging
LOG_LEVEL=info
ENVIRONMENT=development
```

## Testing

```bash
# Esegui i test
cd bff-julia-mobile-api  # o bff-julia-profile-api
make test

# Test con coverage
make test-coverage
```

## Build e Deployment

```bash
# Build del binario
make build

# Build dell'immagine Docker
make docker-build

# Deploy con docker-compose
make docker-run
```

## Monitoraggio

Ogni servizio espone:
- **Health check:** `http://localhost:{port}/health`
- **Metrics (Prometheus):** `http://localhost:{port}/metrics`
- **API Docs (Swagger):** `http://localhost:{port}/swagger/index.html`

## Migrazione da Java

Questi progetti Go sono equivalenti funzionali dei progetti Java Spring Boot originali:
- `bff-julia-mobile-api` (Java) → `bff-julia-mobile-api` (Go)
- `bff-julia-profile-api` (Java) → `bff-julia-profile-api` (Go)

### Principali differenze:
- **Framework:** Spring Boot → Gin
- **DI Container:** Spring → Manual dependency injection
- **ORM:** Spring Data → Azure SDK for Go
- **Logging:** Logback → Zap
- **Build:** Maven → Go modules
- **Performance:** JVM → Native binary (più leggero e veloce)

## License

Comune di Roma
