package auth

import (
	"testing"
)

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
		if hasher.Verify(password, salt, "not-base64-!") {
			t.Error("Verify() should fail for invalid base64 hash")
		}
	})
}
