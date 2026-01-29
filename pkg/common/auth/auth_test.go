package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mocks ---

type mockStore struct {
	users map[string]StoredUser
}

func (m *mockStore) SaveWithPassword(ctx context.Context, data RegisterInput, hash string, salt []byte) error {
	if m.users == nil {
		m.users = make(map[string]StoredUser)
	}
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

type mockSessionManager struct {
	sessions map[string]uuid.UUID
}

func (m *mockSessionManager) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	token := "token-" + userID.String()
	if m.sessions == nil {
		m.sessions = make(map[string]uuid.UUID)
	}
	m.sessions[token] = userID
	return token, time.Now().Add(ttl), nil
}

// --- Tests: Argon2Hasher ---

func TestArgon2Hasher(t *testing.T) {
	hasher := NewArgon2Hasher(DefaultArgonParams)
	password := "my-secure-password"

	t.Run("hash and verify", func(t *testing.T) {
		hash, salt, err := hasher.Hash(password)
		if err != nil {
			t.Fatalf("Hash() failed: %v", err)
		}
		if hash == "" {
			t.Error("Hash() returned empty hash")
		}
		if len(salt) == 0 {
			t.Error("Hash() returned empty salt")
		}

		if !hasher.Verify(password, salt, hash) {
			t.Error("Verify() failed for correct password")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		hash, salt, _ := hasher.Hash(password)
		if hasher.Verify("wrong-password", salt, hash) {
			t.Error("Verify() succeeded for wrong password")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		_, _, err := hasher.Hash("")
		if err == nil {
			t.Error("Hash() should fail for empty password")
		}
	})

	t.Run("short password", func(t *testing.T) {
		_, _, err := hasher.Hash("123")
		if err == nil {
			t.Error("Hash() should fail for short password")
		}
	})

	t.Run("invalid verification input", func(t *testing.T) {
		hash, salt, _ := hasher.Hash(password)
		if hasher.Verify("", salt, hash) {
			t.Error("Verify() should fail for empty password")
		}
		if hasher.Verify(password, nil, hash) {
			t.Error("Verify() should fail for nil salt")
		}
		if hasher.Verify(password, salt, "") {
			t.Error("Verify() should fail for empty hash")
		}
	})
}

// --- Tests: Service ---

func TestAuthService(t *testing.T) {
	store := &mockStore{users: make(map[string]StoredUser)}
	sessions := &mockSessionManager{sessions: make(map[string]uuid.UUID)}

	service, err := NewService(ServiceConfig{
		Store:      store,
		Hasher:     NewArgon2Hasher(DefaultArgonParams),
		Sessions:   sessions,
		SessionTTL: time.Hour,
	})
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	ctx := context.Background()
	input := RegisterInput{
		ID:            uuid.New(),
		Name:          "Test User",
		WalletAddress: "0x123",
		Address:       "Main St",
		Password:      "password123",
	}

	t.Run("Register Success", func(t *testing.T) {
		err := service.Register(ctx, input)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if _, ok := store.users[input.WalletAddress]; !ok {
			t.Error("User not saved in store")
		}
	})

	t.Run("Register Failure - Missing Password", func(t *testing.T) {
		badInput := input
		badInput.Password = ""
		err := service.Register(ctx, badInput)
		if err == nil {
			t.Error("Register should fail with empty password")
		}
	})

	t.Run("Login Success", func(t *testing.T) {
		res, err := service.Login(ctx, input.WalletAddress, input.Password)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if res.Token == "" {
			t.Error("Expected non-empty token")
		}
		if res.User.WalletAddress != input.WalletAddress {
			t.Errorf("Expected wallet address %s, got %s", input.WalletAddress, res.User.WalletAddress)
		}
	})

	t.Run("Login Failure - Wrong Password", func(t *testing.T) {
		_, err := service.Login(ctx, input.WalletAddress, "wrong-password")
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

	t.Run("Validator Usage", func(t *testing.T) {
		validatorCalled := false
		v := func(ctx context.Context, data RegisterInput) error {
			validatorCalled = true
			if data.Name == "Forbidden" {
				return errors.New("forbidden name")
			}
			return nil
		}

		svc, _ := NewService(ServiceConfig{
			Store:     store,
			Sessions:  sessions,
			Validator: v,
		})

		// Test successful validation
		err := svc.Register(ctx, RegisterInput{WalletAddress: "0xabc", Password: "password123", Name: "Allowed"})
		if err != nil {
			t.Errorf("Validator should have allowed registration: %v", err)
		}
		if !validatorCalled {
			t.Error("Validator was not called")
		}

		// Test failed validation
		err = svc.Register(ctx, RegisterInput{WalletAddress: "0xdef", Password: "password123", Name: "Forbidden"})
		if err == nil || err.Error() != "forbidden name" {
			t.Errorf("Validator should have rejected registration, got %v", err)
		}
	})
}

// --- Tests: RedisSessionManager (Logic) ---

func TestRedisSessionManager_NilClient(t *testing.T) {
	var m *RedisSessionManager
	_, _, err := m.Create(context.Background(), uuid.New(), time.Hour)
	if err == nil || err.Error() != "auth: session manager is nil" {
		t.Errorf("Expected nil manager error, got %v", err)
	}

	m = &RedisSessionManager{}
	_, _, err = m.Create(context.Background(), uuid.New(), time.Hour)
	if err == nil || err.Error() != "auth: redis client is not initialized" {
		t.Errorf("Expected uninitialized client error, got %v", err)
	}
}
