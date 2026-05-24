package email

import "testing"

func TestNewEmailAddress_ValidEmail(t *testing.T) {
	addr, err := New("pascal@allen.com")
	if err != nil {
		t.Fatalf("expected no error for valid email, got: %v", err)
	}
	if addr.String() != "pascal@allen.com" {
		t.Fatalf("expected 'pascal@allen.com', got '%s'", addr.String())
	}
}

func TestNewEmailAddress_EmptyString(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty email, got nil")
	}
}

func TestNewEmailAddress_InvalidFormat(t *testing.T) {
	_, err := New("not-an-email")
	if err == nil {
		t.Fatal("expected error for invalid email format, got nil")
	}
}

func TestEmailAddressString(t *testing.T) {
	addr, _ := New("pascal@allen.com")
	if addr.String() != "pascal@allen.com" {
		t.Fatalf("expected 'pascal@allen.com', got '%s'", addr.String())
	}
}
