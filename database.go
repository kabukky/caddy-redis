package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool        *redis.Pool
	redisPrefix = "caddy:"
)

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
	fmt.Println(LogTag, "Setting in database: key:", key, "value:", value)
	reply, err := conn.Do("SET", redisPrefix+key, value)
	if err != nil {
		fmt.Println(LogTag, "Error while setting in Redis:", err)
	}
	fmt.Println(LogTag, "Result:", reply)
}

func get(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()
	fmt.Println(LogTag, "Getting from database: key:", key)
	value, err := redis.Bytes(conn.Do("GET", redisPrefix+key))
	if err != nil {
		return []byte{}, err
	}
	return value, nil
}
