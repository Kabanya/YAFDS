package clients

import "context"

type WalletClient interface {
	GetBalance(ctx context.Context, walletAddress string) (float64, error)
}

type stubWalletClient struct{}

func NewStubWalletClient() WalletClient {
	return &stubWalletClient{}
}

func (c *stubWalletClient) GetBalance(ctx context.Context, walletAddress string) (float64, error) {
	return 1000.0, nil // Stub balance
}
