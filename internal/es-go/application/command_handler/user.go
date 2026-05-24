package command_handler

import (
	"fmt"
	"github.com/pascalallen/es-go/internal/es-go/application/command"
	"github.com/pascalallen/es-go/internal/es-go/domain/email"
	"github.com/pascalallen/es-go/internal/es-go/domain/password"
	"github.com/pascalallen/es-go/internal/es-go/domain/user"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/messaging"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
)

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

	addr, err := email.New(c.EmailAddress)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	hash := password.Create(c.Password)

	u, err := user.Register(c.Id, c.FirstName, c.LastName, addr, hash)
	if err != nil {
		return fmt.Errorf("could not register user: %w", err)
	}

	streamId := fmt.Sprintf("user-%s", c.Id)
	if err := h.EventStore.AppendToStream(streamId, u.Version(), u.UncommittedEvents()); err != nil {
		return fmt.Errorf("could not store events: %w", err)
	}
	u.ClearUncommittedEvents()

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

	addr, err := email.New(c.EmailAddress)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	streamId := fmt.Sprintf("user-%s", c.Id)
	events, err := h.EventStore.ReadFromStream(streamId)
	if err != nil {
		return fmt.Errorf("error reading user stream: %w", err)
	}

	u := user.LoadUserFromEvents(events)

	if err := u.UpdateEmailAddress(addr); err != nil {
		return fmt.Errorf("could not update email address: %w", err)
	}

	if err := h.EventStore.AppendToStream(streamId, u.Version(), u.UncommittedEvents()); err != nil {
		return fmt.Errorf("could not store events: %w", err)
	}
	u.ClearUncommittedEvents()

	return nil
}

type AssignRoleToUserHandler struct {
	EventStore storage.EventStore
}

func (h AssignRoleToUserHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.AssignRoleToUser)
	if !ok {
		return fmt.Errorf("invalid command type passed to AssignRoleToUserHandler: %v", cmd)
	}

	streamId := fmt.Sprintf("user-%s", c.Id)
	events, err := h.EventStore.ReadFromStream(streamId)
	if err != nil {
		return fmt.Errorf("error reading user stream: %w", err)
	}

	u := user.LoadUserFromEvents(events)

	if err := u.AssignRole(c.Role); err != nil {
		return fmt.Errorf("could not assign role: %w", err)
	}

	if err := h.EventStore.AppendToStream(streamId, u.Version(), u.UncommittedEvents()); err != nil {
		return fmt.Errorf("could not store events: %w", err)
	}
	u.ClearUncommittedEvents()

	return nil
}

type DeleteUserHandler struct {
	EventStore storage.EventStore
}

func (h DeleteUserHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.DeleteUser)
	if !ok {
		return fmt.Errorf("invalid command type passed to DeleteUserHandler: %v", cmd)
	}

	streamId := fmt.Sprintf("user-%s", c.Id)
	events, err := h.EventStore.ReadFromStream(streamId)
	if err != nil {
		return fmt.Errorf("error reading user stream: %w", err)
	}

	u := user.LoadUserFromEvents(events)

	if err := u.Delete(); err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	if err := h.EventStore.AppendToStream(streamId, u.Version(), u.UncommittedEvents()); err != nil {
		return fmt.Errorf("could not store events: %w", err)
	}
	u.ClearUncommittedEvents()

	return nil
}
