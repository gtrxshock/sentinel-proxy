package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net"
	"strings"
	"time"
)

type RedisConnector struct {
	requestTimeout time.Duration
}

func NewRedisConnector(requestTimeout int) *RedisConnector {
	return &RedisConnector{requestTimeout: time.Duration(requestTimeout)}
}

func (rc *RedisConnector) Connect(addr string) (net.Conn, error) {
	redisConn, err := net.Dial("tcp", addr)
	if err != nil {
		//log("dial error", err)
		//_ = clientConn.Close()

		return nil, err
	}

	return redisConn, nil
}

func (rc *RedisConnector) Ping(addr string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), rc.requestTimeout*time.Second)
	_, err := redis.NewClient(&redis.Options{Addr: addr}).Ping(ctx).Result()
	cancelFunc()

	if rc.isNoAuthErr(err) {
		return nil
	}

	return err
}

func (rc *RedisConnector) isNoAuthErr(err error) bool {
	if err != nil {
		s := err.Error()
		if strings.HasPrefix(s, "NOAUTH ") {
			return true
		}
	}

	return false
}
