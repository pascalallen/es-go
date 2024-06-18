package query_handler

import (
	"fmt"
	"github.com/pascalallen/es-go/internal/es-go/application/query"
	"github.com/pascalallen/es-go/internal/es-go/domain/user"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/messaging"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
)

type GetUserByIdHandler struct {
	EventStore storage.EventStore
}

func (h GetUserByIdHandler) Handle(qry messaging.Query) (any, error) {
	q, ok := qry.(query.GetUserById)
	if !ok {
		return nil, fmt.Errorf("invalid query type passed to GetUserByIdHandler: %v", qry)
	}

	streamId := fmt.Sprintf("user-%s", q.Id)
	events, err := h.EventStore.ReadFromStream(streamId)
	if err != nil {
		return nil, fmt.Errorf("error attempting to read events from stream: %s", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no events found for user ID: %s", q.Id)
	}

	u := user.LoadUserFromEvents(events)

	return u, nil
}
