package user

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/email"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/domain/password"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"time"
)

type User struct {
	Id                ulid.ULID     `json:"id"`
	FirstName         string        `json:"first_name"`
	LastName          string        `json:"last_name"`
	EmailAddress      string        `json:"email_address"`
	PasswordHash      password.Hash `json:"-"`
	Roles             []role.Role   `json:"roles"`
	CreatedAt         time.Time     `json:"created_at"`
	ModifiedAt        time.Time     `json:"modified_at,omitempty"`
	DeletedAt         time.Time     `json:"deleted_at,omitempty"`
	version           int
	uncommittedEvents []event.Event
}

func Register(id ulid.ULID, firstName, lastName string, addr email.EmailAddress, hash password.Hash) (*User, error) {
	if firstName == "" {
		return nil, fmt.Errorf("first name cannot be empty")
	}
	if lastName == "" {
		return nil, fmt.Errorf("last name cannot be empty")
	}
	u := &User{version: -1}
	u.raise(&event.UserRegistered{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: addr.String(),
		PasswordHash: string(hash),
		OccurredAt:   time.Now(),
	})
	return u, nil
}

func LoadUserFromEvents(events []event.Event) *User {
	u := &User{version: -1}
	for _, evt := range events {
		u.applyEvent(evt)
	}
	return u
}

func (u *User) UpdateEmailAddress(addr email.EmailAddress) error {
	if !u.DeletedAt.IsZero() {
		return fmt.Errorf("cannot update email address of a deleted user")
	}
	u.raise(&event.UserEmailAddressUpdated{
		Id:           u.Id,
		EmailAddress: addr.String(),
		OccurredAt:   time.Now(),
	})
	return nil
}

func (u *User) SetPassword(hash password.Hash) error {
	if !u.DeletedAt.IsZero() {
		return fmt.Errorf("cannot set password of a deleted user")
	}
	u.raise(&event.UserPasswordSet{
		Id:           u.Id,
		PasswordHash: string(hash),
		OccurredAt:   time.Now(),
	})
	return nil
}

func (u *User) AssignRole(r role.Role) error {
	if !u.DeletedAt.IsZero() {
		return fmt.Errorf("cannot assign role to a deleted user")
	}
	for _, existing := range u.Roles {
		if existing.Id == r.Id {
			return fmt.Errorf("role %s is already assigned to this user", r.Name)
		}
	}
	u.raise(&event.UserRoleAssigned{
		Id:         u.Id,
		Role:       r,
		OccurredAt: time.Now(),
	})
	return nil
}

func (u *User) Delete() error {
	if !u.DeletedAt.IsZero() {
		return fmt.Errorf("user is already deleted")
	}
	u.raise(&event.UserDeleted{
		Id:         u.Id,
		OccurredAt: time.Now(),
	})
	return nil
}

func (u *User) Version() int {
	return u.version
}

func (u *User) UncommittedEvents() []event.Event {
	return u.uncommittedEvents
}

func (u *User) ClearUncommittedEvents() {
	u.uncommittedEvents = nil
}

// raise is the write-path helper: mutates state and queues the event for persistence.
// It does NOT increment version — version only changes on the replay path (applyEvent).
func (u *User) raise(evt event.Event) {
	u.applyEventState(evt)
	u.uncommittedEvents = append(u.uncommittedEvents, evt)
}

// applyEvent is the read-path helper: mutates state and increments version.
// Called only by LoadUserFromEvents during event replay.
func (u *User) applyEvent(evt event.Event) {
	u.applyEventState(evt)
	u.version++
}

func (u *User) applyEventState(evt event.Event) {
	switch e := evt.(type) {
	case *event.UserRegistered:
		u.Id = e.Id
		u.FirstName = e.FirstName
		u.LastName = e.LastName
		u.EmailAddress = e.EmailAddress
		u.PasswordHash = password.Hash(e.PasswordHash)
		u.CreatedAt = e.OccurredAt
	case *event.UserEmailAddressUpdated:
		u.EmailAddress = e.EmailAddress
		u.ModifiedAt = e.OccurredAt
	case *event.UserPasswordSet:
		u.PasswordHash = password.Hash(e.PasswordHash)
		u.ModifiedAt = e.OccurredAt
	case *event.UserRoleAssigned:
		u.Roles = append(u.Roles, e.Role)
		u.ModifiedAt = e.OccurredAt
	case *event.UserDeleted:
		u.DeletedAt = e.OccurredAt
	}
}
