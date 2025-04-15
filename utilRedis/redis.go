package utilRedis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	redigo "github.com/gomodule/redigo/redis"
	gookitRedis "github.com/gookit/cache/goredis"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const ErrRedisNil = redis.Nil
const LockerKeyPrefix = "utilRedisLock_"

type RedisClient struct {
	*redis.Client
	ctx                      context.Context
	lockManager              *redsync.Redsync
	maxIdleConn              int
	configAddr               string
	configPassword           string
	configDbNum              int
	configDialTimeoutSeconds int
	configPoolConnections    int
}

var (
	syncLockersCache = map[string]*redsync.Mutex{}
	redisCtx         = context.Background()
)

func NewRedisClient(addr string, password string, db int, dialTimeoutSeconds int, poolConnections int) (rc *RedisClient, err error) {
	if poolConnections < 1 {
		poolConnections = 1
	}
	maxIdleConn := int(poolConnections / 5)

	opts := redis.Options{
		Addr:           addr,
		Password:       password,
		DB:             db,
		DialTimeout:    time.Duration(dialTimeoutSeconds) * time.Second,
		MaxActiveConns: poolConnections,
		MaxIdleConns:   maxIdleConn,
	}
	c := redis.NewClient(&opts)

	err = c.Ping(redisCtx).Err()
	if nil != err {
		return
	}

	rc = &RedisClient{
		Client:                   c,
		ctx:                      redisCtx,
		maxIdleConn:              maxIdleConn,
		configAddr:               addr,
		configPassword:           password,
		configDbNum:              db,
		configDialTimeoutSeconds: dialTimeoutSeconds,
		configPoolConnections:    poolConnections,
	}
	return
}
func (rc *RedisClient) Clone() (*RedisClient, error) {
	return NewRedisClient(rc.configAddr, rc.configPassword, rc.configDbNum, rc.configDialTimeoutSeconds, rc.configPoolConnections)
}

func (rc *RedisClient) Conf() (addr string, password string, db int) {
	return rc.configAddr, rc.configPassword, rc.configDbNum
}

func (rc *RedisClient) CtxDefault() context.Context {
	return rc.ctx
}

func (rc *RedisClient) GetLocker(key string, duration time.Duration, tries int) (locker *redsync.Mutex) {
	key = strings.TrimSpace(key)
	if "" == key {
		return
	}
	lockerName := LockerKeyPrefix + key

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

func (rc *RedisClient) Del(key ...string) (err error) {
	status := rc.Client.Del(rc.ctx, key...)
	return status.Err()
}
func (rc *RedisClient) DelByPattern(pattern string) (err error) {
	script := redis.NewScript(`
local keys = redis.call('KEYS', ARGV[1])
local i 
for i=1,#keys do 
	redis.call('DEL', keys[i])
end
return i
`)
	err = script.Eval(rc.CtxDefault(), rc, []string{}, pattern).Err()
	if ErrorIsNil(err) {
		err = nil
	}
	return
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
func (rc *RedisClient) GetInt(key string) (value int, err error) {
	v, err := rc.Get(key)
	if nil != err {
		return
	}
	value, err = strconv.Atoi(v)
	return
}

func (rc *RedisClient) RedigoConn() (conn redigo.Conn) {
	c, err := rc.Clone()
	if nil != err {
		return
	}
	conn = GoRedisToRedisGoConn(c.Client)
	return
}
func (rc *RedisClient) RedigoPool() (pool *redigo.Pool) {
	pool = &redigo.Pool{Dial: func() (redigo.Conn, error) {
		return rc.RedigoConn(), nil
	}, MaxIdle: rc.maxIdleConn}
	return
}
func (rc *RedisClient) GookitRedis() *gookitRedis.GoRedis {
	return gookitRedis.Connect(rc.Conf())
}

func (rc *RedisClient) BitFill(key string, val int8, start int64, length int64) (err error) {

	script := redis.NewScript(`
local v = 0
if tonumber(ARGV[1]) > 0 then 
	v = 1
end
for i=0,ARGV[3]-1 do
	redis.call('BITFIELD', KEYS[1], 'SET', 'u1', ARGV[2]+i, v)
end
return ARGV[3]
`)
	err = script.Eval(rc.CtxDefault(), rc, []string{key}, val, start, length).Err()

	return
}

func (rc *RedisClient) BitFindSpaceStep(key string, length int64, val int8, start int64, end int64, step int64) (position int64, err error) {
	script := redis.NewScript(`
local p = -1
local v = 0
if tonumber(ARGV[1]) <= 0 then 
	v = 1
end
for i=ARGV[3],ARGV[4]-1,ARGV[5] do
	p = redis.call('BITPOS', KEYS[1], v, i, i+ARGV[2]-1, 'BIT')
	if p < 0 then
		return i
	end
end
return -1
`)
	cmd := script.Eval(rc.CtxDefault(), rc, []string{key}, val, length, start, end, step)
	if nil != cmd.Err() {
		err = cmd.Err()
		return
	}
	position, err = cmd.Int64()
	return
}

func ErrorIsNil(err error) bool {
	return reflect.DeepEqual(err, ErrRedisNil)
}
