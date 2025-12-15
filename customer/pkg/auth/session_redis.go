package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisSessionManager struct {
	client *redis.Client
}

func NewRedisSessionManager(client *redis.Client) *RedisSessionManager {
	return &RedisSessionManager{client: client}
}

func (m *RedisSessionManager) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	if m == nil {
		return "", time.Time{}, errors.New("auth: session manager is nil")
	}
	if m.client == nil {
		return "", time.Time{}, errors.New("auth: redis client is not initialized")
	}
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, fmt.Errorf("auth: failed to generate session token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	key := fmt.Sprintf("session:%s", token)
	expiresAt := time.Now().Add(ttl)
	if err := m.client.Set(ctx, key, userID.String(), ttl).Err(); err != nil {
		return "", time.Time{}, fmt.Errorf("auth: failed to store session token: %w", err)
	}
	return token, expiresAt, nil
}
