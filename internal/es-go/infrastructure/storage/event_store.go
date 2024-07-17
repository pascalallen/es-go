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
	"strings"
)

type EventStore interface {
	AppendToStream(streamId string, event Event) error
	ReadFromStream(streamId string) ([]Event, error)
	RegisterProjection(projection Projection) error
	UnmarshalProjectionResult(name string, result interface{}) error
}

type EventStoreDb struct {
	client           *esdb.Client
	projectionClient *esdb.ProjectionClient
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

	projectionClient, err := esdb.NewProjectionClient(settings)
	if err != nil {
		log.Fatalf("failed to create projection client for event store: %s\n", err)
	}

	return &EventStoreDb{
		client:           client,
		projectionClient: projectionClient,
	}
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

func (s *EventStoreDb) RegisterProjection(projection Projection) error {
	opts := esdb.CreateProjectionOptions{}

	err := s.projectionClient.Create(context.Background(), projection.Name(), projection.Script(), opts)
	if err, ok := esdb.FromError(err); !ok {
		if err.IsErrorCode(esdb.ErrorCodeUnknown) && strings.Contains(err.Err().Error(), "Conflict") {
			log.Printf("projection %s already exists", projection)
		} else {
			return fmt.Errorf("failed to create projection: %s", err)
		}
	}

	return nil
}

func (s *EventStoreDb) UnmarshalProjectionResult(name string, result interface{}) error {
	opts := esdb.GetResultProjectionOptions{}

	value, err := s.projectionClient.GetResult(context.Background(), name, opts)
	if err != nil {
		return fmt.Errorf("failed to get projection result: %s", err)
	}

	jsonContent, err := value.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal projection result: %s", err)
	}

	if err = json.Unmarshal(jsonContent, result); err != nil {
		return fmt.Errorf("failed to unmarshal projection result: %s", err)
	}

	return nil
}
