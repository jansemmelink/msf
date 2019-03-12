package redis

import (
	"time"

	goredis "github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var (
	//ErrTimeout on blocking calls
	ErrTimeout = errors.New("timeout")
)

// Redis ..
type Redis interface {
	BLPOP(queuename string, waitTime int) (string, error)
	LPUSH(queuename string, value string) error
	BRPOP(queuename string, waitTime int) (string, error)
	SET(key string, value string) error
	SETEX(key string, value string, ttl int) error
	SETNX(key string, value string, ttl int) error
	GET(key string) (string, error)
	DEL(key ...string) error
	LLEN(queueName string) (int, error)
	SCAN(key ...interface{}) (int, []string, error)
	Available() int
	MaxActive() int
	Stats() string
}

// NewRedis ..
func NewRedis(network string, server string, poolSize int) (Redis, error) {

	client := goredis.NewClient(&goredis.Options{
		Addr:        server,
		PoolSize:    poolSize,
		ReadTimeout: 10 * time.Second,
	})

	/*p, err := pool.New(network, server, poolSize)

	if err != nil {
		return nil, errors.Wrap(err, "pool creation failed")
	}
	r := redisRadix{pool: p, maxActive: poolSize}

	*/

	r := GoRedis{client: client}

	return r, nil
}
