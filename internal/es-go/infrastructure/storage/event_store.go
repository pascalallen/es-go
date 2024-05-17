package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v4/esdb"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"io"
	"log"
	"os"
)

type EventStore interface {
	AppendToStream(streamId string, event event.Event) error
	ReadFromStream(streamId string, count uint64)
}

type EventStoreDb struct {
	client *esdb.Client
}

func NewEventStoreDb() EventStore {
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

	return &EventStoreDb{client: client}
}

func (s *EventStoreDb) AppendToStream(streamId string, event event.Event) error {
	ropts := esdb.ReadStreamOptions{
		Direction: esdb.Backwards,
		From:      esdb.End{},
	}

	stream, err := s.client.ReadStream(context.Background(), streamId, ropts, 1)
	if err != nil {
		panic(err)
	}

	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx := context.Background()
	opts := esdb.AppendToStreamOptions{
		ExpectedRevision: lastEvent.OriginalStreamRevision(),
	}
	eventData := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   event.EventName(),
		Data:        data,
	}

	_, err = s.client.AppendToStream(ctx, streamId, opts, eventData)
	if err != nil {
		return err
	}

	return nil
}

func (s *EventStoreDb) ReadFromStream(streamId string, count uint64) {
	ctx := context.Background()
	opts := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}

	stream, err := s.client.ReadStream(ctx, streamId, opts, count)
	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		evt, err := stream.Recv()

		if err, ok := esdb.FromError(err); !ok {
			if err.Code() == esdb.ErrorCodeResourceNotFound {
				fmt.Print("Stream not found")
			} else if errors.Is(err, io.EOF) {
				break
			} else {
				panic(err)
			}
		}

		fmt.Printf("Event> %v", evt)
	}
}
