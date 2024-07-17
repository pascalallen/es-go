package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/application/command"
	"github.com/pascalallen/es-go/internal/es-go/application/command_handler"
	"github.com/pascalallen/es-go/internal/es-go/application/projection"
	"github.com/pascalallen/es-go/internal/es-go/application/query"
	"github.com/pascalallen/es-go/internal/es-go/application/query_handler"
	"log"
	"time"
)

func main() {
	container := InitializeContainer()
	defer container.MessageQueueConnection.Close()

	setupProjections(container)

	runConsumers(container)

	tempAppExecution(container)
}

func setupProjections(container Container) {
	eventStore := container.EventStore

	// projection registry
	err := eventStore.RegisterProjection(projection.UserEmailAddresses{})
	if err != nil {
		exitOnError(err)
	}
}

func runConsumers(container Container) {
	commandBus := container.CommandBus
	queryBus := container.QueryBus
	eventStore := container.EventStore

	// command registry
	commandBus.RegisterHandler(command.RegisterUser{}.CommandName(), command_handler.RegisterUserHandler{EventStore: eventStore})
	commandBus.RegisterHandler(command.UpdateUserEmailAddress{}.CommandName(), command_handler.UpdateUserEmailAddressHandler{EventStore: eventStore})

	// query registry
	queryBus.RegisterHandler(query.GetUserById{}.QueryName(), query_handler.GetUserByIdHandler{EventStore: eventStore})

	go commandBus.StartConsuming()
}

func tempAppExecution(container Container) {
	time.Sleep(time.Second * 3)

	userId := ulid.Make()

	go func() {
		// simulate user registration
		registerUserCommand := command.RegisterUser{
			Id:           userId,
			FirstName:    "Pascal",
			LastName:     "Allen",
			EmailAddress: "pascal@allen.com",
		}
		err := container.CommandBus.Execute(registerUserCommand)
		exitOnError(err)

		// simulate email address update
		updateUserEmailCommand := command.UpdateUserEmailAddress{
			Id:           userId,
			EmailAddress: "thomas@allen.com",
		}
		err = container.CommandBus.Execute(updateUserEmailCommand)
		exitOnError(err)
	}()

	time.Sleep(time.Second * 3)

	go func() {
		// simulate querying for user by ID
		getUserByIdQuery := query.GetUserById{Id: userId}
		u, err := container.QueryBus.Fetch(getUserByIdQuery)
		log.Printf("[[[ USER BUILT FROM EVENTS ]]]: %v\n", u)
		exitOnError(err)
	}()

	select {}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
