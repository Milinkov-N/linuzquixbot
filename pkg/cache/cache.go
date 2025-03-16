package cache

import (
	"sync"
	"time"
)

type Cache[T comparable] struct {
	records map[T]*CacheRecord
	timeout time.Duration
	mtx     sync.RWMutex
}

type CacheRecord struct {
	Value  interface{}
	Expiry time.Time
	mtx    sync.RWMutex
}

func NewCache[T comparable](timeout time.Duration) *Cache[T] {
	return &Cache[T]{
		records: map[T]*CacheRecord{},
		timeout: timeout,
	}
}

func NewCacheWithAutoCleanup[T comparable](timeout time.Duration, cleanupInterval time.Duration) *Cache[T] {
	cache := &Cache[T]{
		records: map[T]*CacheRecord{},
		timeout: timeout,
	}

	cache.autoCleanup(cleanupInterval)

	return cache
}

func (c *Cache[T]) Set(key T, value interface{}) {
	record := &CacheRecord{
		Value:  value,
		Expiry: time.Now().Add(c.timeout),
	}

	c.mtx.Lock()
	c.records[key] = record
	c.mtx.Unlock()
}

func (c *Cache[T]) Get(key T) (interface{}, bool) {
	c.mtx.RLock()
	record, exists := c.records[key]

	c.mtx.RUnlock()

	if !exists {
		return nil, false
	}

	record.mtx.RLock()
	defer record.mtx.RUnlock()

	if time.Now().After(record.Expiry) {
		delete(c.records, key)
		return nil, false
	}

	return record.Value, true
}

func (c *Cache[T]) Refresh(key T) {
	c.mtx.RLock()
	record, exists := c.records[key]
	c.mtx.RUnlock()

	if !exists {
		return
	}

	record.mtx.Lock()
	defer record.mtx.Unlock()

	record.Expiry = time.Now().Add(c.timeout)
}

func (c *Cache[T]) autoCleanup(interval time.Duration) {
	tick := time.NewTicker(interval)

	go func() {
		for range tick.C {
			c.mtx.Lock()

			for key, record := range c.records {
				record.mtx.RLock()

				if time.Now().After(record.Expiry) {
					delete(c.records, key)
				}

				record.mtx.RUnlock()
			}

			c.mtx.Unlock()
		}
	}()
}
