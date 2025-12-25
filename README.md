# DAE Services Monorepo

This repository contains the DAE backend services organized as a monorepo.

## Structure

```
.
├── deployments/          # Deployment configurations
│   ├── firebase/        # Firebase emulator configs (Dockerfile, firebase.json)
│   └── observability/   # Observability stack configs (Loki, Tempo, Prometheus, OTEL)
├── libs/                # Shared libraries
│   ├── apperror/        # Application error types
│   ├── prettylog/       # Logging utilities
│   └── utils/           # Common utilities
├── proto/               # Protocol buffer definitions
│   ├── gen/            # Generated protobuf code
│   └── third_party/    # Third-party proto definitions
├── services/            # Microservices
│   └── dae-core/       # Core DAE service
│       ├── cmd/        # Application entrypoints
│       ├── internal/   # Internal packages
│       └── Makefile    # Build and development commands
├── docker-compose.yml   # Infrastructure services (Redis, Firestore, Grafana, etc.)
└── go.work             # Go workspace configuration
```

## Prerequisites

- Go 1.24+
- Protocol Buffers compiler (protoc)
- Docker and Docker Compose

## Development

### Building

From the `services/dae-core` directory:

```bash
make build        # Build the application
make build-linux  # Build for Linux
```

### Running

```bash
make run          # Run the application locally
```

### Testing

```bash
make test              # Run unit tests
make test-integration  # Run integration tests
make cover            # Run tests with coverage
```

### Code Generation

To regenerate protobuf code:

```bash
make gen
```

### Working with the Monorepo

This monorepo uses Go workspaces. The root `go.work` file includes:

- `libs/` - Shared libraries
- `proto/` - Protocol buffer definitions
- `services/dae-core/` - Core service

Dependencies between modules are managed via replace directives in each service's `go.mod`.

## Docker Services

Start all infrastructure services from the repository root:

```bash
docker compose up -d
```

Or use the Makefile shortcuts from the service directory:

```bash
cd services/dae-core
make docker-up
```

This starts:

- Redis (port 6379)
- Firestore Emulator (port 8080)
- Grafana (port 3000)
- Loki, Tempo, Prometheus, OTEL Collector

## Module Structure

### `github.com/deni12345/dae-services/libs`

Shared libraries used across services:

- `apperror`: Application error types and codes
- `prettylog`: Structured logging with OTEL integration
- `utils`: Common utility functions

### `github.com/deni12345/dae-services/proto`

Protocol buffer definitions and generated code for all services.

### `github.com/deni12345/dae-services/services/dae-core`

Main DAE service implementation.

## Migration from dae-core

This repository was migrated from a single-service structure to a monorepo. Key changes:

- `share/` → `libs/`
- `internal/` → `services/dae-core/internal/`
- Proto files moved to root-level `proto/`
- Deployment configs moved to `deployments/`
- Import paths updated to reflect new structure
