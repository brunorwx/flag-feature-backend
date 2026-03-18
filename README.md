# Feature Flag API

## Architecture

The project follows Domain-Driven Design principles with clear separation of concerns:

```
cmd/
├── main.go                                 # Application entry point and HTTP server setup

internal/
├── domain/                                 # Core business logic (entities, interfaces)
│   ├── flag.go                             # FeatureFlag entity with evaluation logic
│   └── flag_repository_interface.go        # Repository interface
├── application/                            # Application services and handlers
│   ├── flag_service.go                     # Business logic orchestration
│   └── flag_handler.go                     # HTTP request/response handling
└── infrastructure/                         # Infrastructure implementations
    └── flag_repository.go                  # In-memory repository implementation
```

### Design Decisions

1. **Domain-Driven Design**: Core business logic (flag evaluation) is in the domain layer, separate from infrastructure
2. **In-Memory Storage**: Easy to swap for a database by implementing the Repository interface.
3. **Thread-Safe Repository**: Uses `sync.RWMutex` for concurrent access.
4. **Chi Router**: Lightweight, fast HTTP router for Go
5. **Repository Pattern**: Abstracts data storage, enabling easy testing and future persistence implementations

## Features

* Create feature flags
* Retrieve flags by key
* Set global on/off state for flags
* Set user-specific overrides
* Rollout percentage-based activation
* Targeted rule-based activation
* Evaluate flags for users
* List all flags

## Feature Flag Model

The `FeatureFlag` entity now includes additional rollout and targeting capabilities:

```json
{
  "key": "string",
  "name": "string",
  "globalEnabled": "bool",
  "rolloutPercentage": "int (0-100)",
  "targetRules": "[][]string",
  "userOverrides": "map[string]bool"
}
```

### Field Description

* **globalEnabled**: Enables/disables the flag globally
* **rolloutPercentage**: Percentage of users (0–100) that should receive the feature
* **targetRules**: List of targeting rules (each rule is a list of conditions, e.g. `[ ["country:PT"], ["plan:premium"] ]`)
* **userOverrides**: Explicit per-user overrides (highest priority)

### Evaluation Priority

1. User override (highest priority)
2. Target rules match
3. Rollout percentage check
4. Global flag state (fallback)

## Quick Start

### 1. Run the Server

```bash
make run
# or
go run ./cmd/main.go
```

The server will start on `http://localhost:8080`

## Project Structure

### Domain Layer (`internal/domain/`)

* **Purpose**: Contains pure business logic independent of infrastructure
* **Contents**: `FeatureFlag` entity with evaluation logic, repository interfaces

### Application Layer (`internal/application/`)

* **Purpose**: Orchestrates domain logic and handles HTTP requests/responses
* **Contents**: `FeatureFlagService` for business operations, HTTP handlers
* **Dependencies**: Domain layer, HTTP router

### Infrastructure Layer (`internal/infrastructure/`)

* **Purpose**: Implements technical concerns (databases, external services)
* **Contents**: In-memory repository implementation with thread safety
* **Dependencies**: Domain layer interfaces

## API Endpoints

### 1. Create a Feature Flag

**POST** `/flags`

```bash
curl -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{
    "key": "new-dashboard",
    "name": "New Dashboard UI",
    "globalEnabled": false,
    "rolloutPercentage": 50,
    "targetRules": [["country:PT"], ["plan:premium"]]
  }'
```

Response (201 Created):

```json
{
  "key": "new-dashboard",
  "name": "New Dashboard UI",
  "globalEnabled": false,
  "rolloutPercentage": 50,
  "targetRules": [["country:PT"], ["plan:premium"]],
  "userOverrides": {}
}
```

### 2. Retrieve a Flag by Key

**GET** `/flags/{key}`

```bash
curl http://localhost:8080/flags/new-dashboard
```

Response (200 OK):

```json
{
  "key": "new-dashboard",
  "name": "New Dashboard UI",
  "globalEnabled": false,
  "rolloutPercentage": 50,
  "targetRules": [["country:PT"], ["plan:premium"]],
  "userOverrides": {}
}
```

### 3. Set Global Flag State

**PUT** `/flags/{key}/global`

```bash
curl -X PUT http://localhost:8080/flags/new-dashboard/global \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### 4. Set User-Specific Override

**PUT** `/flags/{key}/users/{userId}`

```bash
curl -X PUT http://localhost:8080/flags/new-dashboard/users/user123 \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

### 5. Evaluate Flag for a User

**GET** `/evaluate/{key}?userId={userId}`

```bash
curl 'http://localhost:8080/evaluate/new-dashboard?userId=user123'
```

Response (200 OK):

```json
{
  "key": "new-dashboard",
  "userId": "user123",
  "enabled": true
}
```

## Testing curl reqs

```bash
# 1. Create flag with rollout and targeting
curl -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{
    "key": "feature-a",
    "name": "Feature A",
    "globalEnabled": false,
    "rolloutPercentage": 30,
    "targetRules": [["country:PT"]]
  }'

# 2. Enable feature globally
curl -X PUT http://localhost:8080/flags/feature-a/global \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'

# 3. Disable feature for user456
curl -X PUT http://localhost:8080/flags/feature-a/users/user456 \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'

# 4. Evaluate for different users
curl 'http://localhost:8080/evaluate/feature-a?userId=user123'
curl 'http://localhost:8080/evaluate/feature-a?userId=user456'

# 5. List all flags
curl http://localhost:8080/flags
```
