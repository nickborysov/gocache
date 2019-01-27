package cache

import (
	"sync"
	"time"
)

const (
	DefaultCacheCleanDelay = 5 * 60 // min
)

type Cache struct {
	cache           map[string]item
	cacheCleanDelay time.Duration
	lastCleanAt     time.Time
	sync.RWMutex
}

type item struct {
	expiredAt time.Time
	value     interface{}
}

func NewCache() *Cache {
	return &Cache{
		cacheCleanDelay: time.Duration(time.Duration(DefaultCacheCleanDelay) * time.Second),
		cache:           make(map[string]item),
		lastCleanAt:     time.Now(),
	}
}

func (c *Cache) SetCacheCleanDelay(delay int64) *Cache {
	c.cacheCleanDelay = time.Duration(time.Duration(delay) * time.Second)
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	cachedItem, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	if cachedItem.isExpired() {
		return nil, false
	}

	go c.cleanExpiredData()

	return cachedItem.value, true
}

func (c *Cache) GetBool(key string) (bool, bool) {
	value, ok := c.Get(key)
	if !ok {
		return false, false
	}

	intValue, ok := value.(bool)
	if !ok {
		return false, false
	}

	return intValue, true
}

func (c *Cache) GetInt(key string) (int, bool) {
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	intValue, ok := value.(int)
	if !ok {
		return 0, false
	}

	return intValue, true
}

func (c *Cache) GetString(key string) (string, bool) {
	value, ok := c.Get(key)
	if !ok {
		return "", false
	}

	intValue, ok := value.(string)
	if !ok {
		return "", false
	}

	return intValue, true
}

func (c *Cache) SetWithExpiredAt(key string, value interface{}, expiredAt time.Time) {
	c.Lock()
	defer c.Unlock()

	c.cache[key] = item{value: value, expiredAt: expiredAt}

	go c.cleanExpiredData()
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.SetWithExpiredAt(key, value, time.Now().Add(expiration*time.Second))
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.cache, key)

	go c.cleanExpiredData()
}

func (c *Cache) RemoveAll() {
	c.Lock()
	defer c.Unlock()

	c.cache = nil
	c.cache = make(map[string]item)
}

func (c *Cache) cleanExpiredData() {
	if !c.IsNeedToClearCache() {
		// No need clean yet
		return
	}

	c.Lock()
	defer c.Unlock()

	for key, cachedItem := range c.cache {
		if cachedItem.isExpired() {
			delete(c.cache, key)
		}
	}
	c.lastCleanAt = time.Now()
}

func (c *Cache) IsNeedToClearCache() bool {
	c.RLock()
	defer c.RUnlock()
	now := time.Now()
	return c.lastCleanAt.Add(c.cacheCleanDelay * time.Second).Before(now)
}

func (c *Cache) ForceClean() {
	c.cleanExpiredData()
}

func (it *item) isExpired() bool {
	now := time.Now()
	return it.expiredAt.Before(now)
}
