package secrets

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/99designs/keyring"
)

const serviceName = "beeper-cli"

type Credentials struct {
	Name      string    `json:"name,omitempty"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	Get(name string) (*Credentials, error)
	Set(name string, creds Credentials) error
	Delete(name string) error
	List() ([]AccountInfo, error)
}

type AccountInfo struct {
	Name      string
	CreatedAt time.Time
}

type keyringStore struct {
	ring keyring.Keyring
}

type StoreOption func(*storeConfig)

type storeConfig struct {
	fileBackendDir string
}

func WithFileBackend(dir string) StoreOption {
	return func(c *storeConfig) {
		c.fileBackendDir = dir
	}
}

func NewStore(opts ...StoreOption) (Store, error) {
	cfg := &storeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	kc := keyring.Config{
		ServiceName:                    serviceName,
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         false,
		KeychainAccessibleWhenUnlocked: true,
	}

	if cfg.fileBackendDir != "" {
		// File backend is only used for testing
		kc.AllowedBackends = []keyring.BackendType{keyring.FileBackend}
		kc.FileDir = cfg.fileBackendDir
		kc.FilePasswordFunc = func(prompt string) (string, error) {
			return "test-password", nil
		}
	} else {
		// Production: use system keychain
		kc.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
	}

	ring, err := keyring.Open(kc)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &keyringStore{ring: ring}, nil
}

func (s *keyringStore) Get(name string) (*Credentials, error) {
	item, err := s.ring.Get(name)
	if err != nil {
		return nil, fmt.Errorf("credential not found: %s", name)
	}

	var creds Credentials
	if err := json.Unmarshal(item.Data, &creds); err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %w", err)
	}
	creds.Name = name
	return &creds, nil
}

func (s *keyringStore) Set(name string, creds Credentials) error {
	if creds.CreatedAt.IsZero() {
		creds.CreatedAt = time.Now()
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}

	return s.ring.Set(keyring.Item{
		Key:  name,
		Data: data,
	})
}

func (s *keyringStore) Delete(name string) error {
	return s.ring.Remove(name)
}

func (s *keyringStore) List() ([]AccountInfo, error) {
	keys, err := s.ring.Keys()
	if err != nil {
		return nil, err
	}

	accounts := make([]AccountInfo, 0, len(keys))
	for _, key := range keys {
		creds, err := s.Get(key)
		if err != nil {
			continue
		}
		accounts = append(accounts, AccountInfo{
			Name:      key,
			CreatedAt: creds.CreatedAt,
		})
	}
	return accounts, nil
}
