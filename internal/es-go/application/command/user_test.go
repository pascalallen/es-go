package command

import (
	"github.com/oklog/ulid/v2"
	"testing"
)

func TestThatCommandNameReturnsExpectedValueRegisterUser(t *testing.T) {
	cmd := RegisterUser{
		Id:           ulid.Make(),
		FirstName:    "Pascal",
		LastName:     "Allen",
		EmailAddress: "pascal@allen.com",
	}

	if cmd.CommandName() != "RegisterUser" {
		t.Fatal("test assertion failed for RegisterUser.CommandName()")
	}
}

func TestThatCommandNameReturnsExpectedValueRegisterUpdateUserEmailAddress(t *testing.T) {
	cmd := UpdateUserEmailAddress{
		Id:           ulid.Make(),
		EmailAddress: "thomas@allen.com",
	}

	if cmd.CommandName() != "UpdateUserEmailAddress" {
		t.Fatal("test assertion failed for UpdateUserEmailAddress.CommandName()")
	}
}
