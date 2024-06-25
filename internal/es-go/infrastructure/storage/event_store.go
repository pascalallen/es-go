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

type Event interface {
	EventName() string
}

type EventStore interface {
	AppendToStream(streamId string, event Event) error
	ReadFromStream(streamId string) ([]Event, error)
}

type EventStoreDb struct {
	client *esdb.Client
}

func NewEventStoreDb() EventStore {
	connectionString := fmt.Sprintf(
		"esdb://%s:%s?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000",
		os.Getenv("EVENTSTORE_HOST"),
		os.Getenv("EVENTSTORE_HTTP_PORT"),
	)

	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		log.Fatalf("failed to create configuration for event store: %s\n", err)
	}

	client, err := esdb.NewClient(settings)
	if err != nil {
		log.Fatalf("failed to create client for event store: %s\n", err)
	}

	return &EventStoreDb{client: client}
}

func (s *EventStoreDb) AppendToStream(streamId string, event Event) error {
	ropts := esdb.ReadStreamOptions{
		Direction: esdb.Backwards,
		From:      esdb.End{},
	}

	stream, err := s.client.ReadStream(context.Background(), streamId, ropts, 1)
	if err != nil {
		return fmt.Errorf("failed to read from stream for last event: %s", err)
	}

	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err, ok := esdb.FromError(err); !ok {
		if err.Code() == esdb.ErrorCodeResourceNotFound {
			log.Printf("last event stream not found when attempting to append with stream ID: %s", streamId)
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

func (s *EventStoreDb) ReadFromStream(streamId string) ([]Event, error) {
	var events []Event
	position := esdb.Revision(0)

	for {
		opts := esdb.ReadStreamOptions{
			From:      position,
			Direction: esdb.Forwards,
		}

		stream, err := s.client.ReadStream(context.Background(), streamId, opts, 100)
		if err != nil {
			return events, fmt.Errorf("failed to read events from stream: %s", err)
		}

		hasMoreEvents := false

		for {
			evt, err := stream.Recv()
			if err, ok := esdb.FromError(err); !ok {
				if errors.Is(err, io.EOF) {
					break
				}

				if err, ok := esdb.FromError(err); !ok {
					return events, fmt.Errorf("error attempting to stream incoming event: %s", err)
				}
			}

			var e Event
			switch evt.OriginalEvent().EventType {
			case event.UserRegistered{}.EventName():
				e = &event.UserRegistered{}
			case event.UserEmailAddressUpdated{}.EventName():
				e = &event.UserEmailAddressUpdated{}
			default:
				return events, fmt.Errorf("unknown event retrieved: %s", evt.OriginalEvent().EventType)
			}

			err = json.Unmarshal(evt.OriginalEvent().Data, &e)
			if err != nil {
				return events, fmt.Errorf("failed to unmarshal event: %s", err)
			}

			events = append(events, e)

			position = esdb.Revision(evt.OriginalEvent().EventNumber + 1)
			hasMoreEvents = true
		}

		stream.Close()

		if !hasMoreEvents {
			break
		}
	}

	return events, nil
}
