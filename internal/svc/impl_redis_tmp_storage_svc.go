package svc

import (
	"go-backend-template/internal/redis"
	"time"
)

type RedisTmpStorageSvc struct {
	rd *redis.Redis
}

func NewRedisTempStorage(rd *redis.Redis) ITmpStorageSvc {
	return &RedisTmpStorageSvc{
		rd: rd,
	}
}

func (r *RedisTmpStorageSvc) Set(key, value string, expiresAt time.Time) error {
	return r.rd.GetDB().Set(key, value, -time.Now().Sub(expiresAt)).Err()
}

func (r *RedisTmpStorageSvc) Get(key string) (string, error) {
	return r.rd.GetDB().Get(key).Result()
}
