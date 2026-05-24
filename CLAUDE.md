# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Maintenance

When new event types, commands, or infrastructure patterns are added, run `/claude-md-improver` to keep this file current.

## Overview

`es-go` is a reference implementation of Event Sourcing (ES) and Domain-Driven Design (DDD) in Go. The entire development environment runs inside Docker — Go code is executed inside the `go` container with EventStoreDB and RabbitMQ as backing services.

## Commands

All commands run against the Docker environment via the `bin/` scripts:

```bash
bin/up                   # Build images and start all containers (tails logs)
bin/down                 # Stop and remove containers
bin/exec <cmd>           # Run any command inside the go container
```

Run tests (inside the container):
```bash
bin/exec go test ./...
```

Run a single package's tests:
```bash
bin/exec go test github.com/pascalallen/es-go/internal/es-go/application/command/...
```

Run tests with coverage:
```bash
bin/exec go test ./... -covermode=count -coverprofile=coverage.out
bin/exec go tool cover -html=coverage.out -o coverage.html
```

Build the binary:
```bash
bin/exec go build -C cmd/es-go
```

## Architecture

### Layers

```
cmd/es-go/           — entrypoint: wires container, registers commands/queries/projections, starts consumers
internal/es-go/
  domain/            — pure domain types
    email/           — EmailAddress value object (RFC 5322 validation)
    event/           — Event interface + all domain event structs (UserRegistered, UserEmailAddressUpdated,
                       UserPasswordSet, UserRoleAssigned, UserDeleted); every event has OccurredAt time.Time
    password/        — bcrypt Hash value type
    permission/      — Permission value type
    role/            — Role value type
    user/            — User aggregate: Register factory, UpdateEmailAddress, SetPassword, AssignRole, Delete
  application/
    command/         — command structs (RegisterUser, UpdateUserEmailAddress, AssignRoleToUser, DeleteUser)
    command_handler/ — handlers: load aggregate → call method → AppendToStream(streamId, version, events)
    query/           — query structs (GetUserById)
    query_handler/   — handlers: ReadFromStream → LoadUserFromEvents → return aggregate
    projection/      — JS projection scripts run inside EventStoreDB
  infrastructure/
    messaging/       — RabbitMQ command bus (async) and synchronous query bus
    storage/         — EventStoreDB client wrapper (EventStore interface + EventStoreDb impl)
```

### Key patterns

**Aggregate factory + uncommitted events (DDD/ES):**
- `user.Register(...)` is the only way to create a User; it raises a `UserRegistered` event internally
- Aggregate methods (`UpdateEmailAddress`, `SetPassword`, `AssignRole`, `Delete`) call the internal `raise()` helper
- `raise()` mutates state and appends to `uncommittedEvents` but does NOT increment `version`
- `LoadUserFromEvents` replays events via `applyEvent()` which mutates state AND increments `version`
- Command handlers: call aggregate method → `AppendToStream(streamId, u.Version(), u.UncommittedEvents())` → `ClearUncommittedEvents()`

**Optimistic concurrency via aggregate version:**
- `User.version` starts at `-1`. After loading N events: `version == N-1` (EventStoreDB's 0-based revision)
- `AppendToStream` maps `-1` → `esdb.NoStream`, `≥ 0` → `esdb.Revision(version)`
- No read-before-write round trip

**Deterministic event replay:**
- Every event struct carries `OccurredAt time.Time` set at raise-time
- `applyEventState` reads `OccurredAt` for `CreatedAt`/`ModifiedAt`/`DeletedAt` — never calls `time.Now()`

**Dependency direction:**
- `domain/event` defines the `Event` interface
- `infrastructure/storage` imports `domain/event` — never the reverse

### CQRS + Event Sourcing flow

**Commands** (async via RabbitMQ):
1. `CommandBus.Execute` serializes the command and publishes it to the `commands` queue
2. `CommandBus.StartConsuming` deserializes and dispatches to the registered `CommandHandler`
3. The handler validates invariants (e.g., unique email via projection), then calls an aggregate method
4. The aggregate raises events into `uncommittedEvents`
5. The handler calls `EventStore.AppendToStream(streamId, u.Version(), u.UncommittedEvents())`
6. Stream IDs follow the pattern `user-{ulid}`

**Queries** (synchronous):
1. `QueryBus.Fetch` dispatches directly to the registered `QueryHandler`
2. The handler calls `EventStore.ReadFromStream`, then `user.LoadUserFromEvents` to rebuild state

**Projections** (run inside EventStoreDB):
- Written as JavaScript and registered at startup via `EventStore.RegisterProjection`
- Currently used for the `user-email-addresses` unique constraint

### Dependency injection

Uses [Google Wire](https://github.com/google/wire). `wire.go` is the injector (build-tagged `wireinject`); `wire_gen.go` is the generated output. Run `wire` inside the `cmd/es-go` directory to regenerate after modifying providers.

### Adding a new event type

New events require changes in these places:
1. Define the event struct in `internal/es-go/domain/event/user.go` with `EventName()` and `OccurredAt time.Time`
2. Add an `applyEventState` case in `domain/user/user.go`
3. Add a `case` in `EventStoreDb.ReadFromStream` (`infrastructure/storage/event_store.go`)
4. Add the corresponding aggregate method in `domain/user/user.go`
5. Add a `case` in `RabbitMqCommandBus.processCommand` (`infrastructure/messaging/command_bus.go`) if triggered by a new command
6. Register the command handler in `main.go`

### Environment variables

Copy `.env.example` to `.env`. Key variables consumed by the Go app:
- `EVENTSTORE_HOST`, `EVENTSTORE_HTTP_PORT` — EventStoreDB connection
- `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_DEFAULT_USER`, `RABBITMQ_DEFAULT_PASS` — RabbitMQ connection

EventStoreDB UI is available at `http://localhost:2113`; RabbitMQ management UI at `http://localhost:15672`.
