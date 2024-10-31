package cache

import (
	"fmt"
	"sync"
	"time"
)

type TTLKeyedCache[V any] struct {
	items  map[int]item[V]
	mu     sync.Mutex
	nextID int
}

func NewTTLKeyedCache[V any](cleanupInterval time.Duration) *TTLKeyedCache[V] {
	c := &TTLKeyedCache[V]{
		items:  make(map[int]item[V]),
		nextID: 1,
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

func (c *TTLKeyedCache[V]) Set(value V, ttl time.Duration) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Auto increment ID
	id := c.nextID
	c.nextID++

	c.items[id] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}

	return id
}

func (c *TTLKeyedCache[V]) Get(key int) (V, bool) {
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

func (c *TTLKeyedCache[V]) Update(key int, value V, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.items[key]; !found {
		return fmt.Errorf("key %d not found", key)
	}

	c.items[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
	return nil
}

func (c *TTLKeyedCache[V]) Remove(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
