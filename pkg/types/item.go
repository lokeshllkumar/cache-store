package types

import (
	"time"
)

// represents a single cache entry
type CacheItem[K comparable, V any] struct {
	Key        K
	Value      V
	Expiration time.Time
}

// checks if an item has expired
func (i *CacheItem[K, V]) IsExpired() bool {
	return !i.Expiration.IsZero() && time.Now().After(i.Expiration)
}
