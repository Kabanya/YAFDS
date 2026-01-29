package client

import (
	"context"
	"testing"
)

func TestNewStubWalletClient(t *testing.T) {
	client := NewStubWalletClient()
	if client == nil {
		t.Fatal("expected NewStubWalletClient to return a non-nil client")
	}
}

func TestStubWalletClient_GetBalance(t *testing.T) {
	client := NewStubWalletClient()
	ctx := context.Background()
	address := "test-address"

	balance, err := client.GetBalance(ctx, address)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedBalance := 1000.0
	if balance != expectedBalance {
		t.Errorf("expected balance %f, got %f", expectedBalance, balance)
	}
}
