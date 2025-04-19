package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type SessionManager interface {
	Create(ctx context.Context, userID int64, ttl time.Duration) (string, error)
	GetUserID(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
}

type redisSession struct {
	client *redis.Client
}

func NewRedisSession(client *redis.Client) SessionManager {
	return &redisSession{client: client}
}

func (r *redisSession) Create(ctx context.Context, userID int64, ttl time.Duration) (string, error) {
	token := uuid.New().String()
	err := r.client.Set(ctx, "session:"+token, userID, ttl).Err()
	return token, err
}

func (r *redisSession) GetUserID(ctx context.Context, token string) (int64, error) {
	val, err := r.client.Get(ctx, "session:"+token).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

func (r *redisSession) Delete(ctx context.Context, token string) error {
	return r.client.Del(ctx, "session:"+token).Err()
}
