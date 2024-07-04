package utilRedis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type RedisClient struct {
	*redis.Client
	ctx         context.Context
	lockManager *redsync.Redsync
}

func NewRedisClient(addr string, password string, db int, dialTimeoutSeconds int, poolConnections int) (rc *RedisClient, err error) {
	if poolConnections < 1 {
		poolConnections = 1
	}
	maxIdleConns := int(poolConnections / 5)

	opts := redis.Options{
		Addr:           addr,
		Password:       password,
		DB:             db,
		DialTimeout:    time.Duration(dialTimeoutSeconds) * time.Second,
		MaxActiveConns: poolConnections,
		MaxIdleConns:   maxIdleConns,
	}
	c := redis.NewClient(&opts)

	rc = &RedisClient{
		Client: c,
		ctx:    context.Background(),
	}
	return
}

func (rc *RedisClient) GetLocker(key string, duration time.Duration, tries int) (locker *redsync.Mutex) {
	key = strings.TrimSpace(key)
	if "" == key {
		return
	}
	if nil == rc.lockManager {
		rc.lockManager = redsync.New(goredis.NewPool(rc.Client))
	}
	lockerName := "utilRedisLock_" + key
	locker = rc.lockManager.NewMutex(lockerName, redsync.WithExpiry(duration), redsync.WithTries(tries))
	return
}

func (rc *RedisClient) TryLock(key string, duration time.Duration, tries int) error {
	return rc.GetLocker(key, duration, tries).TryLock()
}
func (rc *RedisClient) Lock(key string, duration time.Duration, tries int) error {
	return rc.GetLocker(key, duration, tries).Lock()
}
func (rc *RedisClient) Unlock(key string, duration time.Duration, tries int) (bool, error) {
	return rc.GetLocker(key, duration, tries).Unlock()
}
