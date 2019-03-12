package redis

import (
	"time"

	redigoredis "github.com/gomodule/redigo/redis"
)

//NewPool creates a new redis pool with maxActive set to
// max pool size
func NewPool(server *string, redisMaxActive int) *redigoredis.Pool {
	redis := &redigoredis.Pool{
		MaxActive:   redisMaxActive,
		MaxIdle:     redisMaxActive,
		IdleTimeout: 15 * time.Second,
		Dial: func() (redigoredis.Conn, error) {
			c, err := redigoredis.Dial("tcp", *server,
				redigoredis.DialConnectTimeout(2000*time.Millisecond))
			//redigoredis.DialReadTimeout(20*time.Second))
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redigoredis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redis
}

//NewPoolWithAuth ...
func NewPoolWithAuth(network *string, server *string, redisMaxActive int, password *string) *redigoredis.Pool {
	redis := &redigoredis.Pool{
		MaxActive:   redisMaxActive,
		MaxIdle:     redisMaxActive,
		IdleTimeout: 15 * time.Second,
		Dial: func() (redigoredis.Conn, error) {
			c, err := redigoredis.Dial(*network, *server,
				redigoredis.DialConnectTimeout(2000*time.Millisecond))

			if err != nil {
				return nil, err
			}

			if password != nil && len(*password) != 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redigoredis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redis
}
