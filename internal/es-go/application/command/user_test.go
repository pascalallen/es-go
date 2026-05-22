package command

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"testing"
)

func TestThatCommandNameReturnsExpectedValueRegisterUser(t *testing.T) {
	cmd := RegisterUser{
		Id:           ulid.Make(),
		FirstName:    "Pascal",
		LastName:     "Allen",
		EmailAddress: "pascal@allen.com",
		Password:     "pa$$w0rd",
	}
	if cmd.CommandName() != "RegisterUser" {
		t.Fatal("test assertion failed for RegisterUser.CommandName()")
	}
}

func TestThatCommandNameReturnsExpectedValueUpdateUserEmailAddress(t *testing.T) {
	cmd := UpdateUserEmailAddress{
		Id:           ulid.Make(),
		EmailAddress: "thomas@allen.com",
	}

	if cmd.CommandName() != "UpdateUserEmailAddress" {
		t.Fatal("test assertion failed for UpdateUserEmailAddress.CommandName()")
	}
}

func TestThatCommandNameReturnsExpectedValueAssignRoleToUser(t *testing.T) {
	cmd := AssignRoleToUser{
		Id:   ulid.Make(),
		Role: role.Role{Id: ulid.Make(), Name: "admin"},
	}
	if cmd.CommandName() != "AssignRoleToUser" {
		t.Fatal("test assertion failed for AssignRoleToUser.CommandName()")
	}
}

func TestThatCommandNameReturnsExpectedValueDeleteUser(t *testing.T) {
	cmd := DeleteUser{
		Id: ulid.Make(),
	}
	if cmd.CommandName() != "DeleteUser" {
		t.Fatal("test assertion failed for DeleteUser.CommandName()")
	}
}
