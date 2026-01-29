package client

import (
	"context"
	"testing"
)

func TestNewWalletClient(t *testing.T) {
	client := NewWalletClient()
	if client == nil {
		t.Fatal("expected NewWalletClient to return a non-nil client")
	}
}

func TestWalletClient_GetBalanceWallet(t *testing.T) {
	client := NewWalletClient()
	ctx := context.Background()
	address := "test-address"

	balance, err := client.GetBalanceWallet(ctx, address)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedBalance := 1000.0
	if balance != expectedBalance {
		t.Errorf("expected balance %f, got %f", expectedBalance, balance)
	}
}
