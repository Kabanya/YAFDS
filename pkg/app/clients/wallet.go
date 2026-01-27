package clients

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Kabanya/YAFDS/pkg/common/utils"
)

type WalletClient interface {
	CheckAndDebit(ctx context.Context, walletAddress string, amount float64) (bool, error)
}

type stubWalletClient struct{}

func NewStubWalletClient() WalletClient {
	return &stubWalletClient{}
}

func (c *stubWalletClient) CheckAndDebit(ctx context.Context, walletAddress string, amount float64) (bool, error) {
	logger, _ := utils.Logger()
	logger.Printf("Wallet: checking balance and debiting %f from %s", amount, walletAddress)

	// Simulate wallet service delay
	time.Sleep(10 * time.Millisecond)

	// For demonstration purposes, if wallet address contains "empty", return false (insufficient funds)
	// Otherwise, 95% chance of success
	if walletAddress == "0x_empty" {
		logger.Printf("Wallet: insufficient funds for %s", walletAddress)
		return false, nil
	}

	// Pseudo-random success
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	if r.Float64() < 0.05 {
		logger.Printf("Wallet: random failure for %s", walletAddress)
		return false, fmt.Errorf("wallet service temporary error")
	}

	logger.Printf("Wallet: successfully debited %f from %s", amount, walletAddress)
	return true, nil
}
