package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"customer/pkg/utils"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
}

var DefaultArgonParams = ArgonParams{
	Memory:  utils.Memory64KB,
	Time:    1,
	Threads: utils.NumThreads(6),
	KeyLen:  32,
}

type Argon2Hasher struct {
	params ArgonParams
}

func NewArgon2Hasher(params ArgonParams) *Argon2Hasher {
	if params == (ArgonParams{}) {
		params = DefaultArgonParams
	}
	return &Argon2Hasher{params: params}
}

func (h *Argon2Hasher) WithLogger() *Argon2Hasher {
	return &Argon2Hasher{params: h.params}
}

func (h *Argon2Hasher) Hash(password string) (string, []byte, error) {
	if password == "" {
		return "", nil, errors.New("auth: empty password")
	}
	if len(password) <= 4 {
		return "", nil, errors.New("auth: password must be longer than 4 characters")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", nil, err
	}
	key := argon2.IDKey([]byte(password), salt, h.params.Time, h.params.Memory, h.params.Threads, h.params.KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), salt, nil
}

func (h *Argon2Hasher) Verify(password string, salt []byte, expectedB64 string) bool {
	if password == "" || len(salt) == 0 || expectedB64 == "" {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(expectedB64)
	if err != nil {
		return false
	}
	key := argon2.IDKey([]byte(password), salt, h.params.Time, h.params.Memory, h.params.Threads, uint32(len(expected)))
	return subtle.ConstantTimeCompare(key, expected) == 1
}
