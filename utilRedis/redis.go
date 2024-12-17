package utilRedis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"time"
)

const ErrRedisNil = redis.Nil

type RedisClient struct {
	*redis.Client
	ctx         context.Context
	lockManager *redsync.Redsync
}

var (
	syncLockersCache = map[string]*redsync.Mutex{}
	redisCtx         = context.Background()
)

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

	err = c.Ping(redisCtx).Err()
	if nil != err {
		return
	}
	rc = &RedisClient{
		Client: c,
		ctx:    redisCtx,
	}
	return
}

func (rc *RedisClient) CtxDefault() context.Context {
	return rc.ctx
}

func (rc *RedisClient) GetLocker(key string, duration time.Duration, tries int) (locker *redsync.Mutex) {
	key = strings.TrimSpace(key)
	if "" == key {
		return
	}
	lockerName := "utilRedisLock_" + key

	locker, ok := syncLockersCache[lockerName]
	if !ok {
		if nil == rc.lockManager {
			rc.lockManager = redsync.New(goredis.NewPool(rc.Client))
		}
		locker = rc.lockManager.NewMutex(lockerName)
		syncLockersCache[lockerName] = locker
	}
	if duration > 0 {
		redsync.WithExpiry(duration).Apply(locker)
	}
	if tries > 0 {
		redsync.WithTries(tries).Apply(locker)
	}

	return
}

func (rc *RedisClient) TryLock(key string, duration time.Duration) error {
	return rc.GetLocker(key, duration, 1).TryLock()
}
func (rc *RedisClient) Lock(key string, duration time.Duration, tries int) error {
	return rc.GetLocker(key, duration, tries).Lock()
}
func (rc *RedisClient) Unlock(key string) (bool, error) {
	return rc.GetLocker(key, 0, 1).Unlock()
}
func (rc *RedisClient) LockExtend(key string, duration time.Duration) (bool, error) {
	return rc.GetLocker(key, duration, 1).Extend()
}

func (rc *RedisClient) Set(key string, value interface{}, expirations ...time.Duration) (string, error) {
	expiration := time.Duration(0)
	if len(expirations) > 0 {
		expiration = expirations[0]
	}
	status := rc.Client.Set(rc.ctx, key, value, expiration)
	return status.Result()
}

func (rc *RedisClient) Get(key string) (value string, err error) {
	status := rc.Client.Get(rc.ctx, key)
	return status.Result()
}
func (rc *RedisClient) GetDel(key string) (value string, err error) {
	status := rc.Client.GetDel(rc.ctx, key)
	return status.Result()
}

func (rc *RedisClient) GetString(key string) (value string, err error) {
	value, err = rc.Get(key)

	return
}

func ErrorIsNil(err error) bool {
	return reflect.DeepEqual(err, ErrRedisNil)
}
