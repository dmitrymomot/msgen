package main

import (
	"crypto/tls"
	"time"
	"github.com/gomodule/redigo/redis"
)

func newRedisPool(connStr string, tlsCnf *tls.Config) *redis.Pool {
	return &redis.Pool{
		MaxActive: 10,
		MaxIdle:   10,
		Wait:      true,
		Dial:      setupRedisConnection(connStr, tlsCnf),
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func setupRedisConnection(connStr string, tlsCnf *tls.Config) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		conn, err := redis.DialURL(connStr, redis.DialTLSConfig(tlsCnf))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}