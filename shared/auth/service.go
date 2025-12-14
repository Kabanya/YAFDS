package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const defaultSessionTTL = 30 * time.Minute

// Service coordinates hashing, persistence and sessions.
type Service struct {
	store      Store
	hasher     Hasher
	sessions   SessionManager
	validator  Validator
	logger     Logger
	sessionTTL time.Duration
}

func NewService(cfg ServiceConfig) (*Service, error) {
	if cfg.Store == nil {
		return nil, fmt.Errorf("auth: store is required")
	}
	if cfg.Sessions == nil {
		return nil, fmt.Errorf("auth: session manager is required")
	}
	hasher := cfg.Hasher
	if hasher == nil {
		hasher = NewArgon2Hasher(DefaultArgonParams)
	}
	validator := cfg.Validator
	if validator == nil {
		validator = NoopValidator
	}
	logger := cfg.Logger
	if logger == nil {
		logger = noopLogger{}
	}
	sessTTL := cfg.SessionTTL
	if sessTTL <= 0 {
		sessTTL = defaultSessionTTL
	}

	return &Service{
		store:      cfg.Store,
		hasher:     hasher,
		sessions:   cfg.Sessions,
		validator:  validator,
		logger:     logger,
		sessionTTL: sessTTL,
	}, nil
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	if s == nil {
		return errors.New("auth: service is nil")
	}
	if err := s.ensureInput(input); err != nil {
		return err
	}
	if err := s.validator(ctx, input); err != nil {
		return err
	}
	passwordHash, passwordSalt, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.logger.Printf("auth: hash failed for %s: %v", input.WalletAddress, err)
		return err
	}
	return s.store.SaveWithPassword(ctx, input, passwordHash, passwordSalt)
}

func (s *Service) Login(ctx context.Context, walletAddress string, password string) (LoginResult, error) {
	if s == nil {
		return LoginResult{}, errors.New("auth: service is nil")
	}
	if strings.TrimSpace(walletAddress) == "" {
		return LoginResult{}, errors.New("auth: wallet address is required")
	}
	user, err := s.store.LoadByWalletAddress(ctx, walletAddress)
	if err != nil {
		return LoginResult{}, err
	}
	if !s.hasher.Verify(password, user.PasswordSalt, user.PasswordHash) {
		return LoginResult{}, ErrInvalidCredentials
	}
	token, exp, err := s.sessions.Create(ctx, user.ID, s.sessionTTL)
	if err != nil {
		s.logger.Printf("auth: session create failed for %s: %v", walletAddress, err)
		return LoginResult{}, err
	}
	return LoginResult{User: user, Token: token, Expiration: exp}, nil
}

func (s *Service) ensureInput(input RegisterInput) error {
	if strings.TrimSpace(input.WalletAddress) == "" {
		return errors.New("auth: wallet address is required")
	}
	if strings.TrimSpace(input.Password) == "" {
		return errors.New("auth: password is required")
	}
	return nil
}
