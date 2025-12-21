package id

import (
	"strings"

	"github.com/google/uuid"
)

// FromWallet deterministically derives a UUID from a wallet address so seeded data stays consistent across runs.
func FromWallet(wallet string) uuid.UUID {
	normalized := strings.ToLower(strings.TrimSpace(wallet))
	if normalized == "" {
		return uuid.Nil
	}
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(normalized))
}
