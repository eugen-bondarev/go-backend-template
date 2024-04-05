package svc

import (
	"go-backend-template/internal/redis"
	"time"
)

type RedisTmpStorage struct {
	rd *redis.Redis
}

func NewRedisTempStorage(rd *redis.Redis) ITmpStorage {
	return &RedisTmpStorage{
		rd: rd,
	}
}

func (r *RedisTmpStorage) Set(key, value string, expiresAt time.Time) error {
	return r.rd.GetDB().Set(key, value, -time.Now().Sub(expiresAt)).Err()
}

func (r *RedisTmpStorage) Get(key string) (string, error) {
	return r.rd.GetDB().Get(key).Result()
}
