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

func NewCache(delay int64) *Cache {
	return &Cache{
		cacheCleanDelay: time.Duration(time.Duration(delay) * time.Second),
		cache:           make(map[string]item),
		lastCleanAt:     time.Now(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	cachedItem, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	now := time.Now()
	if cachedItem.expiredAt.Before(now) {
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
	c.SetWithExpiredAt(key, value, time.Now().Add(expiration))
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
	now := time.Now()
	if c.lastCleanAt.Add(c.cacheCleanDelay).After(now) {
		// No need clean yet
		return
	}
	c.Lock()
	defer c.Unlock()
	for key, cachedItem := range c.cache {
		if cachedItem.expiredAt.After(now) {
			delete(c.cache, key)
		}
	}
}
