package policies

import (
	"sync"

	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// FIFO evictor
type FIFOEvictor[K comparable, V any] struct {
	mu    sync.Mutex
	order []K
	items map[K]*types.CacheItem[K, V]
}

// creates a new FIFO evictor
func NewFIFOEvictor[K comparable, V any]() *FIFOEvictor[K, V] {
	return &FIFOEvictor[K, V]{
		order: make([]K, 0),
		items: make(map[K]*types.CacheItem[K, V]),
	}
}

// adds a key to the end of the order if it's new
func (f *FIFOEvictor[K, V]) Update(key K, item *types.CacheItem[K, V]) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.items[key]; !ok {
		f.order = append(f.order, key)
	}
	f.items[key] = item
}

// does nothing for FIFO; empty definition
func (f *FIFOEvictor[K, V]) Access(key K, item *types.CacheItem[K, V]) {}

// deletes a key from the order and items map
func (f *FIFOEvictor[K, V]) Remove(key K) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, k := range f.order {
		if k == key {
			f.order = append(f.order[:i], f.order[i+1:]...)
			delete(f.items, key)
			return
		}
	}
}

// returns the key and item of the first item added
func (f *FIFOEvictor[K, V]) Evict() (K, *types.CacheItem[K, V], bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.order) > 0 {
		key := f.order[0]
		item := f.items[key]
		f.order = f.order[1:]
		delete(f.items, key)
		return key, item, true
	}
	var zero K
	return zero, nil, false
}
