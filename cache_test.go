package main

import (
	"testing"
	"time"

	"github.com/nickborysov/gocache/cache"
)

func TestCache(t *testing.T) {

	var cacheCleanDelay int64 = 1 //sec
	c := cache.NewCache().SetCacheCleanDelay(cacheCleanDelay)

	t.Run("TestCacheStruct", func(t *testing.T) { testCacheStruct(t, c) })
	t.Run("TestCacheInt", func(t *testing.T) { testCacheInt(t, c) })
	t.Run("TestCacheString", func(t *testing.T) { testCacheString(t, c) })
	t.Run("TestCacheBool", func(t *testing.T) { testCacheBool(t, c) })
	t.Run("TestCacheExpiration", func(t *testing.T) { testCacheExpiration(t, c) })
}

func testCacheExpiration(t *testing.T, c *cache.Cache) {
	var key string = "some_value_for_expire"
	var value string = "cache me for expire"
	var expiration time.Duration = 2 // sec
	c.Set(key, value, expiration)

	timeout := time.After(expiration * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-tick:
			cachedValue, ok := c.Get(key)
			if ok && cachedValue == value {
				t.Logf("Value: %+v\n", cachedValue)
			} else {
				t.Errorf("Expected value %+v, current value: %+v\n", value, cachedValue)
			}

		case <-timeout:
			cachedValue, ok := c.Get(key)
			if !ok {
				t.Logf("Value after expiration: %+v\n", cachedValue)
			} else {
				t.Errorf("Expected value %+v, current value: %+v\n", 0, cachedValue)
			}
			return
		}
	}
}

func testCacheStruct(t *testing.T, c *cache.Cache) {
	type SomeStruct struct {
		name string
		age  int
	}

	var key string = "some_value"
	var value SomeStruct = SomeStruct{
		name: "Rocket",
		age:  17,
	}
	var expiration time.Duration = 3 // sec
	c.Set(key, value, expiration)

	{
		cachedValue, ok := c.Get(key)
		if ok && cachedValue == value {
			t.Logf("Value: %+v\n", cachedValue)
		} else {
			t.Errorf("Expected value %+v, current value: %+v\n", value, cachedValue)
		}
	}

	time.Sleep(expiration * time.Second)

	{
		cachedValue, ok := c.Get(key)
		if !ok {
			t.Logf("Value after expiration: %+v\n", cachedValue)
		} else {
			t.Errorf("Expected value %+v, current value: %+v\n", 0, cachedValue)
		}
	}
}

func testCacheInt(t *testing.T, c *cache.Cache) {
	var key string = "some_int_value"
	var value int = 42
	var expiration time.Duration = 3 // sec
	c.Set(key, value, expiration)

	{
		cachedValue, ok := c.GetInt(key)
		if ok && cachedValue == value {
			t.Logf("Value: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", value, cachedValue)
		}
	}

	time.Sleep(expiration * time.Second)

	{
		cachedValue, ok := c.GetInt(key)
		if !ok {
			t.Logf("Value after expiration: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", 0, cachedValue)
		}
	}
}

func testCacheString(t *testing.T, c *cache.Cache) {
	var key string = "some_string_value"
	var value string = "cache me"
	var expiration time.Duration = 3 // sec
	c.Set(key, value, expiration)

	{
		cachedValue, ok := c.GetString(key)
		if ok && cachedValue == value {
			t.Logf("Value: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", value, cachedValue)
		}
	}

	time.Sleep(expiration * time.Second)

	{
		cachedValue, ok := c.GetString(key)
		if !ok {
			t.Logf("Value after expiration: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", 0, cachedValue)
		}
	}
}

func testCacheBool(t *testing.T, c *cache.Cache) {
	var key string = "some_bool_value"
	var value bool = true
	var expiration time.Duration = 3 // sec
	c.Set(key, value, expiration)

	{
		cachedValue, ok := c.GetBool(key)
		if ok && cachedValue == value {
			t.Logf("Value: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", value, cachedValue)
		}
	}

	time.Sleep(expiration * time.Second)

	{
		cachedValue, ok := c.GetBool(key)
		if !ok {
			t.Logf("Value after expiration: %#v\n", cachedValue)
		} else {
			t.Errorf("Expected value %v, current value: %v\n", 0, cachedValue)
		}
	}
}
