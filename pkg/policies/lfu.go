package policies

import (
	"sync"
	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// LFU evictor
type LFUEvictor[K comparable, V any] struct {
	mu sync.Mutex
	count map[K]int
	items map[K]*types.CacheItem[K, V]
}

// creates a new LFU evictor
func NewLFUEvictor[K comparable, V any]() *LFUEvictor[K, V] {
	return &LFUEvictor[K, V]{
		count: make(map[K]int),
		items: make(map[K]*types.CacheItem[K, V]),
	}
}

// sets the count for a new item
func (l *LFUEvictor[K, V]) Update(key K, item *types.CacheItem[K, V]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.count[key]; !ok {
		l.count[key] = 1
	}
	l.items[key] = item
}

// increments the access count for a key
func (l *LFUEvictor[K, V]) Access(key K, item *types.CacheItem[K, V]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.count[key]++
}

// deletes the item and count for a key
func (l *LFUEvictor[K, V]) Remove(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.count, key)
	delete(l.items, key)
}

// finds and returns key and item with the lowest access count
func (l *LFUEvictor[K, V]) Evict() (K, *types.CacheItem[K, V], bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.count) == 0 {
		var zero K
		return zero, nil, false
	}

	var minKey K
	minCount := -1
	for k, c := range l.count {
		if minCount == -1 || c < minCount {
			minCount = c
			minKey = k
		}
	}

	if item, ok := l.items[minKey]; ok {
		delete(l.count, minKey)
		delete(l.items, minKey)
		return minKey, item, true
	}
	var zero K
	return zero, nil, false
}
