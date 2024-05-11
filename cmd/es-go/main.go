package main

import (
	"errors"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/domain/user"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
	"io"
)

func main() {
	eventStore := storage.NewEventStore()

	u := &user.User{}

	// simulate u registration
	userId := ulid.Make()
	firstName := "Pascal"
	lastName := "Allen"
	emailAddress := "pascal@allen.com"
	registerEvent := event.UserRegistered{
		Id:           userId,
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: emailAddress,
	}
	u.ApplyEvent(registerEvent)
	streamId := fmt.Sprintf("user-%s", userId)
	err := eventStore.AppendToStream(streamId, registerEvent)
	if err != nil {
		panic(err)
	}

	// simulate email address update
	newEmailAddress := "thomas@allen.com"
	emailUpdateEvent := event.UserEmailAddressUpdated{
		EmailAddress: newEmailAddress,
	}
	u.ApplyEvent(emailUpdateEvent)
	err = eventStore.AppendToStream(streamId, emailUpdateEvent)
	if err != nil {
		panic(err)
	}

	stream, _ := eventStore.GetStream(streamId, 10)
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
