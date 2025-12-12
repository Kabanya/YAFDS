package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"runtime"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
}

var defaultParams = ArgonParams{
	Memory:  64 * 1024, // 64 MB
	Time:    1,
	Threads: uint8(runtime.NumCPU()),
	KeyLen:  32,
}

func HashPassword(password string, params ArgonParams) (hashB64 string, salt []byte, err error) {
	if password == "" {
		return "", nil, errors.New("empty password")
	}
	if len(password) <= 4 {
		return "", nil, errors.New("less 4 symbols in passwords")
	}
	salt = make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return "", nil, err
	}
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)
	return base64.RawStdEncoding.EncodeToString(key), salt, nil
}

func VerifyPassword(password string, params ArgonParams, salt []byte, expectedB64 string) bool {
	if password == "" || len(password) == 4 || len(salt) == 0 || expectedB64 == "" {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(expectedB64)
	if err != nil {
		return false
	}
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, uint32(len(expected)))
	return subtle.ConstantTimeCompare(key, expected) == 1
}
