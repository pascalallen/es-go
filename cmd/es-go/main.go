package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pascalallen/es-go/internal/es-go/application/event"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
	"io"
)

func main() {
	testEvent := event.TestEvent{
		Id:            uuid.NewString(),
		ImportantData: "I wrote my first event!",
	}

	eventStore := storage.NewEventStore()

	streamId := fmt.Sprintf("my-stream-%s", uuid.New().String())

	err := eventStore.AppendToStream(streamId, testEvent)
	if err != nil {
		panic(err)
	}

	stream, err := eventStore.GetStream(streamId, 10)
	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		evt, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(evt.OriginalEvent())
	}
}
