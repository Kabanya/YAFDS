package id

import (
	"testing"

	"github.com/google/uuid"
)

func TestFromWallet(t *testing.T) {
	tests := []struct {
		name   string
		wallet string
		want   uuid.UUID
	}{
		{
			name:   "empty wallet",
			wallet: "",
			want:   uuid.Nil,
		},
		{
			name:   "whitespace wallet",
			wallet: "   ",
			want:   uuid.Nil,
		},
		{
			name:   "valid wallet",
			wallet: "0x1234567890abcdef1234567890abcdef12345678",
			want:   uuid.NewSHA1(uuid.NameSpaceURL, []byte("0x1234567890abcdef1234567890abcdef12345678")),
		},
		{
			name:   "case insensitive",
			wallet: "0xABCDEF",
			want:   uuid.NewSHA1(uuid.NameSpaceURL, []byte("0xabcdef")),
		},
		{
			name:   "trimmed whitespace",
			wallet: "  0xabcdef  ",
			want:   uuid.NewSHA1(uuid.NameSpaceURL, []byte("0xabcdef")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromWallet(tt.wallet); got != tt.want {
				t.Errorf("FromWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}
