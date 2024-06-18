package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/application/command"
	"github.com/pascalallen/es-go/internal/es-go/application/command_handler"
	"github.com/pascalallen/es-go/internal/es-go/application/query"
	"github.com/pascalallen/es-go/internal/es-go/application/query_handler"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/messaging"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
	"log"
	"time"
)

func main() {
	eventStore, err := storage.NewEventStoreDb()
	exitOnError(err)
	rabbitMqConn, err := messaging.NewRabbitMQConnection()
	exitOnError(err)
	defer rabbitMqConn.Close()
	commandBus, err := messaging.NewRabbitMqCommandBus(rabbitMqConn)
	exitOnError(err)
	queryBus := messaging.NewSynchronousQueryBus()

	// command registry
	commandBus.RegisterHandler(command.RegisterUser{}.CommandName(), command_handler.RegisterUserHandler{EventStore: eventStore})
	commandBus.RegisterHandler(command.UpdateUserEmailAddress{}.CommandName(), command_handler.UpdateUserEmailAddressHandler{EventStore: eventStore})

	go commandBus.StartConsuming()

	// query registry
	queryBus.RegisterHandler(query.GetUserById{}.QueryName(), query_handler.GetUserByIdHandler{EventStore: eventStore})

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
		err = commandBus.Execute(registerUserCommand)
		exitOnError(err)

		// simulate email address update
		updateUserEmailCommand := command.UpdateUserEmailAddress{
			Id:           userId,
			EmailAddress: "thomas@allen.com",
		}
		err = commandBus.Execute(updateUserEmailCommand)
		exitOnError(err)
	}()

	time.Sleep(time.Second * 3)

	go func() {
		// simulate querying for user by ID
		getUserByIdQuery := query.GetUserById{Id: userId}
		u, err := queryBus.Fetch(getUserByIdQuery)
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
