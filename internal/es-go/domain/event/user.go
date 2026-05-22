package event

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"time"
)

type UserRegistered struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
	PasswordHash string    `json:"password_hash"`
	OccurredAt   time.Time `json:"occurred_at"`
}

func (e UserRegistered) EventName() string {
	return "UserRegistered"
}

type UserEmailAddressUpdated struct {
	Id           ulid.ULID `json:"id"`
	EmailAddress string    `json:"email_address"`
	OccurredAt   time.Time `json:"occurred_at"`
}

func (e UserEmailAddressUpdated) EventName() string {
	return "UserEmailAddressUpdated"
}

type UserPasswordSet struct {
	Id           ulid.ULID `json:"id"`
	PasswordHash string    `json:"password_hash"`
	OccurredAt   time.Time `json:"occurred_at"`
}

func (e UserPasswordSet) EventName() string {
	return "UserPasswordSet"
}

type UserRoleAssigned struct {
	Id         ulid.ULID `json:"id"`
	Role       role.Role `json:"role"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e UserRoleAssigned) EventName() string {
	return "UserRoleAssigned"
}

type UserDeleted struct {
	Id         ulid.ULID `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (e UserDeleted) EventName() string {
	return "UserDeleted"
}
