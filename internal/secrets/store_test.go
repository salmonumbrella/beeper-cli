package secrets

import (
	"testing"
	"time"
)

func TestCredentialsRoundtrip(t *testing.T) {
	// Use file backend for testing (no keyring access needed)
	store, err := NewStore(WithFileBackend(t.TempDir()))
	if err != nil {
		t.Fatalf("NewStore() error: %v", err)
	}

	creds := Credentials{
		Token:     "test-token-123",
		CreatedAt: time.Now().Truncate(time.Second),
	}

	if err := store.Set("test", creds); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	got, err := store.Get("test")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}

	if got.Token != creds.Token {
		t.Errorf("Token = %q, want %q", got.Token, creds.Token)
	}
}

func TestStoreList(t *testing.T) {
	store, err := NewStore(WithFileBackend(t.TempDir()))
	if err != nil {
		t.Fatalf("NewStore() error: %v", err)
	}

	_ = store.Set("alpha", Credentials{Token: "a"})
	_ = store.Set("beta", Credentials{Token: "b"})

	accounts, err := store.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("List() returned %d accounts, want 2", len(accounts))
	}
}

func TestStoreDelete(t *testing.T) {
	store, err := NewStore(WithFileBackend(t.TempDir()))
	if err != nil {
		t.Fatalf("NewStore() error: %v", err)
	}

	_ = store.Set("test", Credentials{Token: "t"})
	if err := store.Delete("test"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err = store.Get("test")
	if err == nil {
		t.Error("Get() after Delete() should return error")
	}
}
