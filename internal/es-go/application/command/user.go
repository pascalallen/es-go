package command

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
)

type RegisterUser struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
	Password     string    `json:"password"`
}

func (c RegisterUser) CommandName() string {
	return "RegisterUser"
}

type UpdateUserEmailAddress struct {
	Id           ulid.ULID `json:"id"`
	EmailAddress string    `json:"email_address"`
}

func (c UpdateUserEmailAddress) CommandName() string {
	return "UpdateUserEmailAddress"
}

type AssignRoleToUser struct {
	Id   ulid.ULID `json:"id"`
	Role role.Role `json:"role"`
}

func (c AssignRoleToUser) CommandName() string {
	return "AssignRoleToUser"
}

type DeleteUser struct {
	Id ulid.ULID `json:"id"`
}

func (c DeleteUser) CommandName() string {
	return "DeleteUser"
}
