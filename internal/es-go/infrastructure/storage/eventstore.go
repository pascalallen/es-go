package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"log"
	"os"
)

type EventStore struct {
	client *esdb.Client
}

func NewEventStore() EventStore {
	connectionString := fmt.Sprintf(
		"esdb://%s:%s?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000",
		"eventstore",
		os.Getenv("EVENTSTORE_HTTP_PORT"),
	)

	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("failed to create configuration for event store: %s", err)
	}

	client, err := esdb.NewClient(settings)
	if err != nil {
		log.Fatalf("failed to create client for event store: %s", err)
	}

	return EventStore{client: client}
}

func (s EventStore) AppendToStream(streamId string, event event.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx := context.Background()
	opts := esdb.AppendToStreamOptions{}
	eventData := esdb.EventData{
		ContentType: esdb.JsonContentType,
		EventType:   event.EventName(),
		Data:        data,
	}

	_, err = s.client.AppendToStream(ctx, streamId, opts, eventData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: potentially abstract ReadStream return type
func (s EventStore) GetStream(streamId string, count int) (*esdb.ReadStream, error) {
	ctx := context.Background()
	opts := esdb.ReadStreamOptions{}

	stream, err := s.client.ReadStream(ctx, streamId, opts, uint64(count))
	if err != nil {
		return nil, err
	}

	return stream, nil
}
