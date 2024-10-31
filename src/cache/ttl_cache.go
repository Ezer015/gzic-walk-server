package cache

import (
	"sync"
	"time"
)

type TTLCache[K comparable, V any] struct {
	items map[K]item[V]
	mu    sync.Mutex
}

func NewTTLCache[K comparable, V any](cleanupInterval time.Duration) *TTLCache[K, V] {
	c := &TTLCache[K, V]{
		items: make(map[K]item[V]),
	}

	go func() {
		for range time.Tick(cleanupInterval) {
			c.mu.Lock()

			// Iterate over the cache items and delete expired ones.
			for key, item := range c.items {
				if item.isExpired() {
					delete(c.items, key)
				}
			}

			c.mu.Unlock()
		}
	}()

	return c
}

func (c *TTLCache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return item.value, false
	}
	if item.isExpired() {
		delete(c.items, key)
		return item.value, false
	}

	return item.value, true
}

func (c *TTLCache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
