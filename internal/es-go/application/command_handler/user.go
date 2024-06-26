package command_handler

import (
	"fmt"
	"github.com/pascalallen/es-go/internal/es-go/application/command"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/messaging"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
)

type ProjectionState struct {
	EmailAddresses map[string]string `json:"email_addresses"`
}

type RegisterUserHandler struct {
	EventStore storage.EventStore
}

func (h RegisterUserHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.RegisterUser)
	if !ok {
		return fmt.Errorf("invalid command type passed to RegisterUserHandler: %v", cmd)
	}

	var result ProjectionState
	if err := h.EventStore.UnmarshalProjectionResult("user-email-addresses", &result); err != nil {
		return fmt.Errorf("error getting projection result: %v", err)
	}

	for emailAddress := range result.EmailAddresses {
		if emailAddress == c.EmailAddress {
			return fmt.Errorf("email address %s is already registered", c.EmailAddress)
		}
	}

	registerEvent := event.UserRegistered{
		Id:           c.Id,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		EmailAddress: c.EmailAddress,
	}
	streamId := fmt.Sprintf("user-%s", c.Id)
	err := h.EventStore.AppendToStream(streamId, registerEvent)
	if err != nil {
		return fmt.Errorf("could not store UserRegistered event: %w", err)
	}

	return nil
}

type UpdateUserEmailAddressHandler struct {
	EventStore storage.EventStore
}

func (h UpdateUserEmailAddressHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.UpdateUserEmailAddress)
	if !ok {
		return fmt.Errorf("invalid command type passed to UpdateUserEmailAddressHandler: %v", cmd)
	}

	var result ProjectionState
	if err := h.EventStore.UnmarshalProjectionResult("user-email-addresses", &result); err != nil {
		return fmt.Errorf("error getting projection result: %v", err)
	}

	for emailAddress := range result.EmailAddresses {
		if emailAddress == c.EmailAddress {
			return fmt.Errorf("could not update user. email address %s is already taken", c.EmailAddress)
		}
	}

	emailUpdateEvent := event.UserEmailAddressUpdated{
		Id:           c.Id,
		EmailAddress: c.EmailAddress,
	}
	streamId := fmt.Sprintf("user-%s", c.Id)
	err := h.EventStore.AppendToStream(streamId, emailUpdateEvent)
	if err != nil {
		return fmt.Errorf("could not store UserEmailAddressUpdated event: %w", err)
	}

	return nil
}
