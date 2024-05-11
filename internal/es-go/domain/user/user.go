package user

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/domain/password"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"time"
)

type User struct {
	Id           ulid.ULID             `json:"id"`
	FirstName    string                `json:"first_name"`
	LastName     string                `json:"last_name"`
	EmailAddress string                `json:"email_address"`
	PasswordHash password.PasswordHash `json:"-"`
	Roles        []role.Role           `json:"roles"`
	CreatedAt    time.Time             `json:"created_at"`
	ModifiedAt   time.Time             `json:"modified_at,omitempty"`
	DeletedAt    time.Time             `json:"deleted_at,omitempty"`
}

func (u *User) ApplyEvent(evt event.Event) {
	switch e := evt.(type) {
	case event.UserRegistered:
		u.Id = e.Id
		u.FirstName = e.FirstName
		u.LastName = e.LastName
		u.EmailAddress = e.EmailAddress
		u.CreatedAt = time.Now()
	case event.UserEmailAddressUpdated:
		u.EmailAddress = e.EmailAddress
		u.ModifiedAt = time.Now()
	}
}
