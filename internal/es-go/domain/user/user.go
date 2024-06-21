package user

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/domain/password"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"github.com/pascalallen/es-go/internal/es-go/infrastructure/storage"
	"time"
)

type User struct {
	Id           ulid.ULID     `json:"id"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	EmailAddress string        `json:"email_address"`
	PasswordHash password.Hash `json:"-"`
	Roles        []role.Role   `json:"roles"`
	CreatedAt    time.Time     `json:"created_at"`
	ModifiedAt   time.Time     `json:"modified_at,omitempty"`
	DeletedAt    time.Time     `json:"deleted_at,omitempty"`
}

func LoadUserFromEvents(events []storage.Event) *User {
	u := &User{}

	for _, evt := range events {
		u.applyEvent(evt)
	}

	return u
}

func (u *User) applyEvent(evt storage.Event) {
	switch evt.EventName() {
	case event.UserRegistered{}.EventName():
		e := evt.(*event.UserRegistered)
		u.Id = e.Id
		u.FirstName = e.FirstName
		u.LastName = e.LastName
		u.EmailAddress = e.EmailAddress
		u.CreatedAt = time.Now()
	case event.UserEmailAddressUpdated{}.EventName():
		e := evt.(*event.UserEmailAddressUpdated)
		u.Id = e.Id
		u.EmailAddress = e.EmailAddress
		u.ModifiedAt = time.Now()
	}
}
