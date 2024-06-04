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
	ReadFromStream(streamId string, count uint64) error
}

type EventStoreDb struct {
	client *esdb.Client
}

func NewEventStoreDb() (EventStore, error) {
	connectionString := fmt.Sprintf(
		"esdb://%s:%s?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000",
		"eventstore",
		os.Getenv("EVENTSTORE_HTTP_PORT"),
	)

	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration for event store: %s", err)
	}

	client, err := esdb.NewClient(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for event store: %s", err)
	}

	return &EventStoreDb{client: client}, nil
}

func (s *EventStoreDb) AppendToStream(streamId string, event event.Event) error {
	ropts := esdb.ReadStreamOptions{
		Direction: esdb.Backwards,
		From:      esdb.End{},
	}

	stream, err := s.client.ReadStream(context.Background(), streamId, ropts, 1)
	if err != nil {
		return fmt.Errorf("failed to read from stream: %s", err)
	}

	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err, ok := esdb.FromError(err); !ok {
		if err.Code() == esdb.ErrorCodeResourceNotFound {
			log.Printf("event stream not found when appending to stream with ID: %s", streamId)
		} else {
			return fmt.Errorf("failed to get last event from stream: %s", err)
		}
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event for stream: %s", err)
	}

	opts := esdb.AppendToStreamOptions{}
	if lastEvent != nil {
		opts = esdb.AppendToStreamOptions{
			ExpectedRevision: lastEvent.OriginalStreamRevision(),
		}
	}
	eventData := esdb.EventData{
		ContentType: esdb.ContentTypeJson,
		EventType:   event.EventName(),
		Data:        data,
	}

	_, err = s.client.AppendToStream(context.Background(), streamId, opts, eventData)
	if err != nil {
		return fmt.Errorf("failed to append event to stream: %s", err)
	}

	return nil
}

// work in progress
func (s *EventStoreDb) ReadFromStream(streamId string, count uint64) error {
	ctx := context.Background()
	opts := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}

	stream, err := s.client.ReadStream(ctx, streamId, opts, count)
	if err != nil {
		return fmt.Errorf("failed to read from stream: %s", err)
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

	return nil
}
