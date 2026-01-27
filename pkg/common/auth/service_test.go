package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockStore struct {
	users map[string]StoredUser
}

func (m *mockStore) SaveWithPassword(ctx context.Context, data RegisterInput, hash string, salt []byte) error {
	m.users[data.WalletAddress] = StoredUser{
		ID:            data.ID,
		Name:          data.Name,
		WalletAddress: data.WalletAddress,
		Address:       data.Address,
		PasswordHash:  hash,
		PasswordSalt:  salt,
	}
	return nil
}

func (m *mockStore) LoadByWalletAddress(ctx context.Context, walletAddress string) (StoredUser, error) {
	u, ok := m.users[walletAddress]
	if !ok {
		return StoredUser{}, errors.New("not found")
	}
	return u, nil
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, []byte, error) {
	return "hashed-" + password, []byte("salt"), nil
}
func (m *mockHasher) Verify(password string, salt []byte, expected string) bool {
	return expected == "hashed-"+password
}

type mockSessions struct{}

func (m *mockSessions) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	return "token-" + userID.String(), time.Now().Add(ttl), nil
}

func TestAuthService(t *testing.T) {
	store := &mockStore{users: make(map[string]StoredUser)}
	hasher := &mockHasher{}
	sessions := &mockSessions{}

	service, _ := NewService(ServiceConfig{
		Store:      store,
		Hasher:     hasher,
		Sessions:   sessions,
		SessionTTL: time.Hour,
	})

	ctx := context.Background()
	input := RegisterInput{
		ID:            uuid.New(),
		Name:          "Test User",
		WalletAddress: "0x123",
		Address:       "Main St",
		Password:      "password123",
	}

	t.Run("Register", func(t *testing.T) {
		err := service.Register(ctx, input)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if _, ok := store.users[input.WalletAddress]; !ok {
			t.Error("User not saved in store")
		}
	})

	t.Run("Login Success", func(t *testing.T) {
		res, err := service.Login(ctx, input.WalletAddress, input.Password)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if res.Token != "token-"+input.ID.String() {
			t.Errorf("Expected token-..., got %s", res.Token)
		}
	})

	t.Run("Login Failure - Wrong Password", func(t *testing.T) {
		_, err := service.Login(ctx, input.WalletAddress, "wrong")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("Expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("Login Failure - Not Found", func(t *testing.T) {
		_, err := service.Login(ctx, "nonexistent", "password")
		if err == nil {
			t.Error("Expected error for nonexistent user")
		}
	})
}
