package cache

import (
	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// interface for different eviction strategies
type EvictionPolicy[K comparable, V any] interface {
	Update(key K, item *types.CacheItem[K, V])
	Access(key K, item *types.CacheItem[K, V])
	Remove(key K)
	Evict() (K, *types.CacheItem[K, V], bool)
}
