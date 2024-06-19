package utilCache

import (
	"github.com/gookit/cache"
	"github.com/gookit/cache/gcache"
	"github.com/gookit/cache/goredis"
	"time"
)

type Cache struct {
	cacheManager   *cache.Manager
	cacheKeyPrefix string
	expire         time.Duration
}

func NewCache(keyPrefix string, expire time.Duration) *Cache {
	return &Cache{cacheManager: cache.NewManager(), cacheKeyPrefix: keyPrefix, expire: expire}
}

func (c *Cache) RegisterStoreFile(dir string) *Cache {
	c.cacheManager.UnregisterAll()
	c.cacheManager.Register(cache.DvrFile, cache.NewFileCache(dir, c.cacheKeyPrefix))
	c.cacheManager.DefaultUse(cache.DvrFile)
	return c
}

func (c *Cache) RegisterStoreMemory(size int) *Cache {
	c.cacheManager.UnregisterAll()
	c.cacheManager.Register(gcache.Name, gcache.New(size))
	c.cacheManager.DefaultUse(gcache.Name)

	return c
}
func (c *Cache) RegisterStoreRedis(addr, pwd string, dbNum int) *Cache {
	c.cacheManager.UnregisterAll()
	goRedis := goredis.Connect(addr, pwd, dbNum)
	goRedis.WithOptions(cache.WithPrefix(c.cacheKeyPrefix))
	c.cacheManager.Register(goredis.Name, goRedis)
	c.cacheManager.DefaultUse(goredis.Name)
	return c
}

func (c *Cache) SetExpire(ttl time.Duration) *Cache {
	c.expire = ttl
	return c
}

func (c *Cache) Close() {
	c.cacheManager.Close()
}

// Has cache key
func (c *Cache) Has(key string) bool {
	return c.cacheManager.Has(key)
}

// Get value by key
func (c *Cache) Get(key string) interface{} {
	return c.cacheManager.Get(key)
}
func (c *Cache) GetBool(key string) (v bool) {
	vt := c.cacheManager.Get(key)
	if nil == vt {
		return
	}

	v, ok := vt.(bool)
	if !ok {
		v = false
	}

	return
}
func (c *Cache) GetString(key string) (v string, ok bool) {
	vt := c.cacheManager.Get(key)
	if nil == vt {
		return
	}
	v, ok = vt.(string)
	return
}
func (c *Cache) GetInt(key string) (v int, ok bool) {
	vt := c.cacheManager.Get(key)
	if nil == vt {
		return
	}
	v, ok = vt.(int)
	return
}
func (c *Cache) GetInt64(key string) (v int64, ok bool) {
	vt := c.cacheManager.Get(key)
	if nil == vt {
		return
	}
	v, ok = vt.(int64)
	return
}

// Set value by key
func (c *Cache) Set(key string, val interface{}, ttl ...time.Duration) error {
	expire := c.expire
	if len(ttl) > 0 {
		expire = ttl[0]
	}
	return c.cacheManager.Set(key, val, expire)
}

// Del value by key
func (c *Cache) Del(key string) error {
	return c.cacheManager.Del(key)
}

// GetMulti values by keys
func (c *Cache) GetMulti(keys []string) map[string]interface{} {
	return c.cacheManager.GetMulti(keys)
}

// SetMulti values
func (c *Cache) SetMulti(mv map[string]interface{}, ttl ...time.Duration) error {
	expire := c.expire
	if len(ttl) > 0 {
		expire = ttl[0]
	}
	return c.cacheManager.SetMulti(mv, expire)
}

// DelMulti values by keys
func (c *Cache) DelMulti(keys []string) error {
	return c.cacheManager.DelMulti(keys)
}

func (c *Cache) ClearAll() {
	c.cacheManager.ClearAll()
}
