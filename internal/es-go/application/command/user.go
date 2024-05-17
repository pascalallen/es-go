package command

import "github.com/oklog/ulid/v2"

type RegisterUser struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
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
