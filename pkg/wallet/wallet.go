package wallet

import "context"

type WalletClient interface {
	GetBalanceWallet(ctx context.Context, walletAddress string) (float64, error)
	PayToWallet(ctx context.Context, walletAddress string, amount float64) error
}

type walletClient struct{}

func NewWalletClient() WalletClient {
	return &walletClient{}
}

func (c *walletClient) GetBalanceWallet(ctx context.Context, walletAddress string) (float64, error) {
	return 1000.0, nil
}

func (c *walletClient) PayToWallet(ctx context.Context, walletAddress string, amount float64) error {
	return nil
}
