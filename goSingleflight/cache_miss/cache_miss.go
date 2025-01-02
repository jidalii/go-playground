package cache_miss

import (
	"context"
	"sync"
	"time"
	// "log"

	redis "github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	once   sync.Once
}

func (c *Cache) SetupRedis() {
	c.once.Do(func() {
		c.Client = redis.NewClient(&redis.Options{
			Addr:            "localhost:6379",
			Password:        "",
			DB:              0,
			ReadTimeout:     time.Second * 5,
			WriteTimeout:    time.Second * 5,
			ConnMaxIdleTime: time.Minute * 5,
			PoolSize:        40,
			MinIdleConns:    10,
			MaxRetries:      3,
		})

		_, err := c.Client.Ping(context.Background()).Result()
		if err != nil {
			panic(err)
		}
		// log.Debug("Redis connected")
	})
}

func (c *Cache) GetUser(id string) (string, error) {
	// simulate redis read delay
	user, err := c.Client.Get(context.Background(), id).Result()
	if err != nil {
		return "", err
	}

	return user, nil
}

func (c *Cache) SetUser(id string, user string) error {
	err := c.Client.Set(context.Background(), id, user, 0).Err()
	if err != nil {
		return err
	}
	return nil

}

func MockDBUserQuery(id string) (string, error) {
	// simulate complex db read delay
	time.Sleep(time.Millisecond * 200)
	return "Alice", nil
}
