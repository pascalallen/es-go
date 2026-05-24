Add a new domain event to the es-go codebase. This requires changes in exactly 5 places — work through them in order and run `bin/exec go build -C cmd/es-go` at the end to verify.

If the user hasn't specified the event name and fields, ask before proceeding.

## Checklist

1. **Define the event struct** in `internal/es-go/domain/event/user.go`
   - Add a struct with the event's fields (use `ulid.ULID` for IDs)
   - Implement `EventName() string` returning a PascalCase name matching the struct name

2. **Register the event in `ReadFromStream`** in `internal/es-go/infrastructure/storage/event_store.go`
   - Add a `case` to the type switch that instantiates the new event struct

3. **Register the command in `processCommand`** in `internal/es-go/infrastructure/messaging/command_bus.go`
   - Add a `case` to the type switch for the corresponding command (skip if this event has no new command)

4. **Wire up the handler in `main.go`** (`cmd/es-go/main.go`)
   - Add `commandBus.RegisterHandler(command.YourCommand{}.CommandName(), command_handler.YourHandler{EventStore: eventStore})`
   - (Skip if reusing an existing command)

5. **Apply the event in the aggregate** in `internal/es-go/domain/user/user.go`
   - Add a `case` to `applyEvent` that updates the `User` struct fields from the event

## Verification

```bash
bin/exec go build -C cmd/es-go
bin/exec go test ./...
```

If the build fails, a missing `case` in one of the type switches is the most likely cause — check steps 2 and 3 first.
