package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	rd "github.com/go-redis/redis"
)

type Redis struct {
	db *rd.Client
}

func NewRedis(host, port, password string) (Redis, error) {
	db := rd.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	return Redis{
		db: db,
	}, db.ClientID().Err()
}

func (r *Redis) GetDB() *rd.Client {
	return r.db
}
