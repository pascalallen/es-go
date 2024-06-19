package messaging

import "fmt"

type Query interface {
	QueryName() string
}

type QueryHandler interface {
	Handle(query Query) (any, error)
}

type QueryBus interface {
	RegisterHandler(queryType string, handler QueryHandler)
	Fetch(qry Query) (any, error)
}

type SynchronousQueryBus struct {
	handlers map[string]QueryHandler
}

func NewSynchronousQueryBus() QueryBus {
	return &SynchronousQueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

func (bus *SynchronousQueryBus) RegisterHandler(queryType string, handler QueryHandler) {
	bus.handlers[queryType] = handler
}

func (bus *SynchronousQueryBus) Fetch(query Query) (any, error) {
	handler, found := bus.handlers[query.QueryName()]
	if !found {
		return nil, fmt.Errorf("no handler registered for query type: %s", query.QueryName())
	}

	results, err := handler.Handle(query)
	if err != nil {
		return nil, fmt.Errorf("error calling query handler: %s", err)
	}

	return results, nil
}
