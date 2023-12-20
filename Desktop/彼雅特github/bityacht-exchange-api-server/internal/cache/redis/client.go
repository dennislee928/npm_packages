package redis

import (
	"bityacht-exchange-api-server/configs"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var lock sync.RWMutex

const Nil = redis.Nil

func Client() *redis.Client {
	lock.RLock()
	if client != nil {
		defer lock.RUnlock()

		return client
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()

	if client != nil {
		return client
	}

	client = Connect()
	return client
}

// Connect to Redis by config of this project
func Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", configs.Config.Cache.Redis.Host, configs.Config.Cache.Redis.Port),
		Username:     configs.Config.Cache.Redis.Username,
		Password:     configs.Config.Cache.Redis.Password,
		DB:           configs.Config.Cache.Redis.DB,
		ReadTimeout:  configs.Config.Cache.Redis.ReadTimeout,
		WriteTimeout: configs.Config.Cache.Redis.WriteTimeout,
		PoolSize:     configs.Config.Cache.Redis.MaxConnections,
	})
}
