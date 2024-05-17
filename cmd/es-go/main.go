package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/application/command"
	"github.com/pascalallen/es-go/internal/es-go/application/command_handler"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/messaging"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
)

func main() {
	eventStore := storage.NewEventStoreDb()
	rabbitMqConn := messaging.NewRabbitMQConnection()
	defer rabbitMqConn.Close()
	commandBus := messaging.NewRabbitMqCommandBus(rabbitMqConn)

	// command registry
	commandBus.RegisterHandler(command.RegisterUser{}.CommandName(), command_handler.RegisterUserHandler{EventStore: eventStore})
	commandBus.RegisterHandler(command.UpdateUserEmailAddress{}.CommandName(), command_handler.UpdateUserEmailAddressHandler{EventStore: eventStore})

	go commandBus.StartConsuming()

	// simulate user registration
	userId := ulid.Make()
	registerUserCommand := command.RegisterUser{
		Id:           userId,
		FirstName:    "Pascal",
		LastName:     "Allen",
		EmailAddress: "pascal@allen.com",
	}
	commandBus.Execute(registerUserCommand)

	// simulate email address update
	updateUserEmailCommand := command.UpdateUserEmailAddress{
		Id:           userId,
		EmailAddress: "thomas@allen.com",
	}
	commandBus.Execute(updateUserEmailCommand)
}
