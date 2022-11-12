package redis

import (
	"auth/domain"
	"github.com/go-redis/redis"
	"time"
)

type redisRepo struct {
	Client *redis.Client
}

func NewRedisRepo(cl *redis.Client) domain.JwtTokenRepo {
	return &redisRepo{Client: cl}
}

func (r *redisRepo) InsertTokenRepo(key, token string, ttl time.Duration) error {
	if err := r.Client.Set(key, token, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisRepo) FindTokenRepo(key, token string) (bool, error) {
	value, err := r.Client.Get(key).Result()
	if err != nil {
		return false, err
	}
	return value == token, nil
}
