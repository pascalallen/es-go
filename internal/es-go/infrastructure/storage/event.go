package storage

import "github.com/pascalallen/es-go/internal/es-go/domain/event"

// Event is a temporary alias for domain/event.Event.
// It will be removed in the next task once all callers use event.Event directly.
type Event = event.Event
