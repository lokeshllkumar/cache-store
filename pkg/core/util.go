package cache

import (
	"time"

	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// applies the eviction policy to remove an item
func (c *Cache[K, V]) evict() {
	removalKey, removalItem, ok := c.evictor.Evict()
	if ok {
		c.items.Delete(removalKey)
		c.metrics.IncrementEvictions()
		if c.onEvictHook != nil {
			c.onEvictHook(removalKey, removalItem.Value)
		}
	}
}

// returns the number of items in the cache
func (c *Cache[K, V]) size() int {
	ct := 0
	// callback function in Range to iterate through the c.items map
	c.items.Range(func(_, _ interface{}) bool {
		ct++
		return true
	})
	return ct
}

// starts a background goroutine to clean up expired items
func (c *Cache[K, V]) startCleaner(interval time.Duration) {
	c.cleanupWG.Add(1)
	go func() {
		defer c.cleanupWG.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanExpired()
			case <-c.cleanupStop:
				return
			}
		}
	}()
}

// iterates through cache and removes expired items

func (c *Cache[K, V]) cleanExpired() {
	c.items.Range(func(key, value interface{}) bool {
		item := value.(*types.CacheItem[K, V])
		if item.IsExpired() {
			c.Delete(key.(K))
			if c.onEvictHook != nil {
				c.onEvictHook(key.(K), item.Value)
			}
		}
		return true
	})
}
