package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

//GoRedis is a wrapper around a REDIS Client
type GoRedis struct {
	client *redis.Client
}

//Available ...
func (r GoRedis) Available() int {
	return int(r.client.PoolStats().IdleConns)
}

//Stats ...
func (r GoRedis) Stats() string {

	stats := r.client.PoolStats()

	return fmt.Sprintf("%+v", stats)
}

// BLPOP ...
func (r GoRedis) BLPOP(queuename string, waitTime int) (string, error) {

	d, err := time.ParseDuration(fmt.Sprintf("%ds", waitTime))

	if err != nil {
		return "", errors.Wrap(err, "BLPOP parseDuration failed")
	}

	result := r.client.BLPop(d, queuename)

	if result.Err() == redis.Nil {
		return "", nil
	} else if result.Err() != nil {
		return "", errors.Wrap(result.Err(), "failed sending BLPOP")
	}

	if len(result.Val()) == 1 {
		return "", nil
	}

	return result.Val()[1], nil

}

//BRPOP ...
func (r GoRedis) BRPOP(queuename string, waitTime int) (string, error) {

	d, err := time.ParseDuration(fmt.Sprintf("%ds", waitTime))

	if err != nil {
		return "", errors.Wrap(err, "BRPOP parseDuration failed")
	}

	result := r.client.BRPop(d, queuename)

	if result.Err() == redis.Nil {
		return "", nil
	} else if result.Err() != nil {
		return "", errors.Wrap(result.Err(), "failed sending BRPOP")
	}

	if len(result.Val()) == 1 {
		return "", nil
	}

	return result.Val()[1], nil
}

// LPUSH ...
func (r GoRedis) LPUSH(queuename string, value string) error {

	intCmd := r.client.LPush(queuename, value)

	return errors.Wrap(intCmd.Err(), "Failed to LPUSH")

}

// SET ..
func (r GoRedis) SET(key string, value string) error {

	return errors.Wrap(r.client.Set(key, value, 0*time.Second).Err(), "Failed to SET value")

}

// SETEX ..
func (r GoRedis) SETEX(key string, value string, ttl int) error {

	return errors.Wrap(r.client.Set(key, value, time.Duration(ttl)*time.Second).Err(), "Failed to SET value")

}

// SETNX ..
func (r GoRedis) SETNX(key string, value string, ttl int) error {
	return errors.Wrap(r.client.SetNX(key, value, time.Duration(ttl)*time.Second).Err(), "Failed to SETNX value")
}

// GET ..
func (r GoRedis) GET(key string) (string, error) {

	res := r.client.Get(key)

	if res.Err() != nil {
		return "", errors.Wrap(res.Err(), "Failed to SET")
	}

	return res.Val(), nil

}

//DEL ..
func (r GoRedis) DEL(key ...string) error {

	cmd := r.client.Del(key...)

	return errors.Wrap(cmd.Err(), "Failed to DEL")
}

//LLEN ...
func (r GoRedis) LLEN(queueName string) (int, error) {

	cmd := r.client.LLen(queueName)

	return int(cmd.Val()), cmd.Err()

}

//SCAN ...
func (r GoRedis) SCAN(args ...interface{}) (cursor int, values []string, err error) {

	return 0, make([]string, 0), nil

	/*resp := r.pool.Cmd("SCAN", args)
	cmd := r.client.SScan()

	if resp.Err != nil {
		return 0, nil, errors.Wrap(err, "failed to run scan")
	}

	if resp.IsType(redis.Array) {
		arr, err := resp.Array()

		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to get array")
		}

		cursor, err = arr[0].Int()

		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to convert cursor to int")
		}

		if arr[1].IsType(redis.Array) {

			v, err := arr[1].Array()
			if err != nil {
				return 0, nil, errors.Wrap(err, "failed to convert values to array")
			}
			values = make([]string, len(v))
			for index, value := range v {
				values[index], err = value.Str()
				if err != nil {
					return 0, nil, errors.Wrap(err, "failed to convert values to string")
				}
			}

		} else {
			return 0, nil, errors.New("cannot convert values to array")
		}

		return cursor, values, nil
	}

	return 0, nil, nil*/
}

//MaxActive ...
func (r GoRedis) MaxActive() int {
	return int(r.client.PoolStats().TotalConns)
}
