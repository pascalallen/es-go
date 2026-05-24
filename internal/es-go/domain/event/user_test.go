package event

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"testing"
	"time"
)

func TestThatEventNameReturnsExpectedValueUserRegistered(t *testing.T) {
	e := UserRegistered{
		Id:           ulid.Make(),
		FirstName:    "Pascal",
		LastName:     "Allen",
		EmailAddress: "pascal@allen.com",
		PasswordHash: "hash",
		OccurredAt:   time.Now(),
	}
	if e.EventName() != "UserRegistered" {
		t.Fatal("test assertion failed for UserRegistered.EventName()")
	}
}

func TestThatEventNameReturnsExpectedValueUserEmailAddressUpdated(t *testing.T) {
	e := UserEmailAddressUpdated{
		Id:           ulid.Make(),
		EmailAddress: "pascal@allen.com",
		OccurredAt:   time.Now(),
	}
	if e.EventName() != "UserEmailAddressUpdated" {
		t.Fatal("test assertion failed for UserEmailAddressUpdated.EventName()")
	}
}

func TestThatEventNameReturnsExpectedValueUserPasswordSet(t *testing.T) {
	e := UserPasswordSet{
		Id:           ulid.Make(),
		PasswordHash: "hash",
		OccurredAt:   time.Now(),
	}
	if e.EventName() != "UserPasswordSet" {
		t.Fatal("test assertion failed for UserPasswordSet.EventName()")
	}
}

func TestThatEventNameReturnsExpectedValueUserRoleAssigned(t *testing.T) {
	e := UserRoleAssigned{
		Id:         ulid.Make(),
		Role:       role.Role{Id: ulid.Make(), Name: "admin"},
		OccurredAt: time.Now(),
	}
	if e.EventName() != "UserRoleAssigned" {
		t.Fatal("test assertion failed for UserRoleAssigned.EventName()")
	}
}

func TestThatEventNameReturnsExpectedValueUserDeleted(t *testing.T) {
	e := UserDeleted{
		Id:         ulid.Make(),
		OccurredAt: time.Now(),
	}
	if e.EventName() != "UserDeleted" {
		t.Fatal("test assertion failed for UserDeleted.EventName()")
	}
}
