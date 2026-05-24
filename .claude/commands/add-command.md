Add a new CQRS command to the es-go codebase. This requires changes in 4 places — work through them in order.

If the user hasn't specified the command name and fields, ask before proceeding.

## Checklist

1. **Define the command struct** in `internal/es-go/application/command/user.go`
   - Add a struct with the command's fields (use `ulid.ULID` for IDs)
   - Implement `CommandName() string` returning a descriptive name (e.g. `"RegisterUser"`)

2. **Register the command in `processCommand`** in `internal/es-go/infrastructure/messaging/command_bus.go`
   - Add a `case` to the type switch that instantiates the new command struct

3. **Create the command handler** in `internal/es-go/application/command_handler/user.go`
   - Implement the `CommandHandler` interface: `Handle(cmd messaging.Command) error`
   - Type-assert the incoming command, check projection state for any invariants, then append the resulting event to EventStoreDB via `h.EventStore.AppendToStream`
   - Stream ID convention: `"user-" + id.String()`
   - If this command produces a new event type, follow `/new-event` first

4. **Register the handler in `main.go`** (`cmd/es-go/main.go`)
   - Add `commandBus.RegisterHandler(command.YourCommand{}.CommandName(), command_handler.YourHandler{EventStore: eventStore})`

## Verification

```bash
bin/exec go build -C cmd/es-go
bin/exec go test ./...
```

If the build fails, a missing `case` in the `processCommand` type switch (step 2) is the most likely cause.
