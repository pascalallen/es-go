package user

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/es-go/internal/es-go/domain/email"
	"github.com/pascalallen/es-go/internal/es-go/domain/event"
	"github.com/pascalallen/es-go/internal/es-go/domain/password"
	"github.com/pascalallen/es-go/internal/es-go/domain/role"
	"testing"
	"time"
)

func TestUserRegister_HappyPath(t *testing.T) {
	id := ulid.Make()
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")

	u, err := Register(id, "Pascal", "Allen", addr, hash)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if u.Id != id {
		t.Fatalf("expected Id %v, got %v", id, u.Id)
	}
	if u.FirstName != "Pascal" {
		t.Fatalf("expected FirstName 'Pascal', got '%s'", u.FirstName)
	}
	if u.LastName != "Allen" {
		t.Fatalf("expected LastName 'Allen', got '%s'", u.LastName)
	}
	if u.EmailAddress != "pascal@allen.com" {
		t.Fatalf("expected EmailAddress 'pascal@allen.com', got '%s'", u.EmailAddress)
	}
	if len(u.UncommittedEvents()) != 1 {
		t.Fatalf("expected 1 uncommitted event, got %d", len(u.UncommittedEvents()))
	}
	if _, ok := u.UncommittedEvents()[0].(*event.UserRegistered); !ok {
		t.Fatal("expected first uncommitted event to be *event.UserRegistered")
	}
	if u.Version() != -1 {
		t.Fatalf("expected version -1 after Register (nothing persisted), got %d", u.Version())
	}
}

func TestUserRegister_EmptyFirstName(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")

	_, err := Register(ulid.Make(), "", "Allen", addr, hash)
	if err == nil {
		t.Fatal("expected error for empty first name, got nil")
	}
}

func TestUserRegister_EmptyLastName(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")

	_, err := Register(ulid.Make(), "Pascal", "", addr, hash)
	if err == nil {
		t.Fatal("expected error for empty last name, got nil")
	}
}

func TestUserUpdateEmailAddress_HappyPath(t *testing.T) {
	id := ulid.Make()
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(id, "Pascal", "Allen", addr, hash)
	u.ClearUncommittedEvents()

	newAddr, _ := email.New("thomas@allen.com")
	err := u.UpdateEmailAddress(newAddr)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if u.EmailAddress != "thomas@allen.com" {
		t.Fatalf("expected 'thomas@allen.com', got '%s'", u.EmailAddress)
	}
	if len(u.UncommittedEvents()) != 1 {
		t.Fatalf("expected 1 uncommitted event, got %d", len(u.UncommittedEvents()))
	}
	if _, ok := u.UncommittedEvents()[0].(*event.UserEmailAddressUpdated); !ok {
		t.Fatal("expected uncommitted event to be *event.UserEmailAddressUpdated")
	}
}

func TestUserUpdateEmailAddress_AlreadyDeleted(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	_ = u.Delete()

	newAddr, _ := email.New("thomas@allen.com")
	err := u.UpdateEmailAddress(newAddr)
	if err == nil {
		t.Fatal("expected error when updating deleted user, got nil")
	}
}

func TestUserSetPassword_HappyPath(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	u.ClearUncommittedEvents()

	newHash := password.Create("newpa$$w0rd")
	err := u.SetPassword(newHash)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(u.UncommittedEvents()) != 1 {
		t.Fatalf("expected 1 uncommitted event, got %d", len(u.UncommittedEvents()))
	}
	if _, ok := u.UncommittedEvents()[0].(*event.UserPasswordSet); !ok {
		t.Fatal("expected uncommitted event to be *event.UserPasswordSet")
	}
}

func TestUserAssignRole_HappyPath(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	u.ClearUncommittedEvents()

	r := role.Role{Id: ulid.Make(), Name: "admin"}
	err := u.AssignRole(r)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(u.Roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(u.Roles))
	}
	if len(u.UncommittedEvents()) != 1 {
		t.Fatalf("expected 1 uncommitted event, got %d", len(u.UncommittedEvents()))
	}
	if _, ok := u.UncommittedEvents()[0].(*event.UserRoleAssigned); !ok {
		t.Fatal("expected uncommitted event to be *event.UserRoleAssigned")
	}
}

func TestUserAssignRole_DuplicateRoleId(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	r := role.Role{Id: ulid.Make(), Name: "admin"}
	_ = u.AssignRole(r)
	u.ClearUncommittedEvents()

	err := u.AssignRole(r)
	if err == nil {
		t.Fatal("expected error for duplicate role ID, got nil")
	}
}

func TestUserDelete_HappyPath(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	u.ClearUncommittedEvents()

	err := u.Delete()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if u.DeletedAt.IsZero() {
		t.Fatal("expected DeletedAt to be set, got zero time")
	}
	if len(u.UncommittedEvents()) != 1 {
		t.Fatalf("expected 1 uncommitted event, got %d", len(u.UncommittedEvents()))
	}
	if _, ok := u.UncommittedEvents()[0].(*event.UserDeleted); !ok {
		t.Fatal("expected uncommitted event to be *event.UserDeleted")
	}
}

func TestUserDelete_AlreadyDeleted(t *testing.T) {
	addr, _ := email.New("pascal@allen.com")
	hash := password.Create("pa$$w0rd")
	u, _ := Register(ulid.Make(), "Pascal", "Allen", addr, hash)
	_ = u.Delete()

	err := u.Delete()
	if err == nil {
		t.Fatal("expected error when deleting already-deleted user, got nil")
	}
}

func TestLoadUserFromEvents_DeterministicTimestamps(t *testing.T) {
	id := ulid.Make()
	registeredAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	events := []event.Event{
		&event.UserRegistered{
			Id:           id,
			FirstName:    "Pascal",
			LastName:     "Allen",
			EmailAddress: "pascal@allen.com",
			PasswordHash: "hash",
			OccurredAt:   registeredAt,
		},
		&event.UserEmailAddressUpdated{
			Id:           id,
			EmailAddress: "thomas@allen.com",
			OccurredAt:   updatedAt,
		},
	}

	u := LoadUserFromEvents(events)

	if !u.CreatedAt.Equal(registeredAt) {
		t.Fatalf("expected CreatedAt %v, got %v", registeredAt, u.CreatedAt)
	}
	if !u.ModifiedAt.Equal(updatedAt) {
		t.Fatalf("expected ModifiedAt %v, got %v", updatedAt, u.ModifiedAt)
	}
	if len(u.UncommittedEvents()) != 0 {
		t.Fatalf("expected no uncommitted events after LoadUserFromEvents, got %d", len(u.UncommittedEvents()))
	}
	if u.Version() != 1 {
		t.Fatalf("expected version 1 after 2 events, got %d", u.Version())
	}
}
