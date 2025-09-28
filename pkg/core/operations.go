package cache

import (
	"time"

	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// set adds a new item to the cache or update an existing one
func (c *Cache[K, V]) Set(key K, value V, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &types.CacheItem[K, V]{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(duration),
	}
	c.evictor.Update(key, item)
	c.items.Store(key, item)

	// eviction check
	if c.size() > c.capacity {
		c.evict()
	}
}

// retrieves an item from the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.items.Load(key); ok {
		item := val.(*types.CacheItem[K, V])
		if !item.IsExpired() {
			c.evictor.Access(key, item)
			c.metrics.IncrementHits()
			return item.Value, true
		}
	}
	c.metrics.IncrementMisses()
	var zero V
	return zero, false
}

// removes an item from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items.Delete(key)
	c.evictor.Remove(key)
}

// gracefully shuts down cache cleaner
func (c *Cache[K, V]) Stop() {
	if c.cleanupEnabled {
		close(c.cleanupStop)
		c.cleanupWG.Wait()
	}
}

// regsiters a callback function for eviction events
func (c *Cache[K, V]) OnEvict(hook func(K, V)) {
	c.onEvictHook = hook
}

// return current cache metrics (getter)
func (c *Cache[K, V]) Metrics() Metrics {
	return *c.metrics
}
