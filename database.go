package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func newPool(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func set(key string, value string) {
	conn := pool.Get()
	defer conn.Close()
	fmt.Println(LOG_TAG, "Setting in database: key:", key, "value:", value)
	reply, err := conn.Do("SET", "caddy:"+key, value)
	if err != nil {
		fmt.Println(LOG_TAG, "Error while setting in Redis:", err)
	}
	fmt.Println(LOG_TAG, "Result:", reply)
}
