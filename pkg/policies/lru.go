package policies

import (
	"container/list"
	"sync"
	"github.com/lokeshllkumar/cache-store/v1/pkg/types"
)

// LRU evictor
type LRUEvictor[K comparable, V any] struct {
	mu   sync.Mutex
	list *list.List
	keys map[K]*list.Element
}

// creates a new LRU evictor
func NewLRUEvictor[K comparable, V any]() *LRUEvictor[K, V] {
	return &LRUEvictor[K, V]{
		list: list.New(),
		keys: make(map[K]*list.Element),
	}
}

// adds a new item or moves an existing item to the front of the list
func (l *LRUEvictor[K, V]) Update(key K, item *types.CacheItem[K, V]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element, ok := l.keys[key]; ok {
		l.list.MoveToFront(element)
		element.Value.(*types.CacheItem[K, V]).Value = item.Value
		element.Value.(*types.CacheItem[K, V]).Expiration = item.Expiration
	} else {
		element := l.list.PushFront(item)
		l.keys[key] = element
	}
}

// moves an accessed item to the front of the list
func (l *LRUEvictor[K, V]) Access(key K, item *types.CacheItem[K, V]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element, ok := l.keys[key]; ok {
		l.list.MoveToFront(element)
	}
}

// deletes an item from the linked list and map
func (l *LRUEvictor[K, V]) Remove(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element, ok := l.keys[key]; ok {
		l.list.Remove(element)
		delete(l.keys, key)
	}
}

// returns the keys and item of the least recently used item from the back of the list
func (l *LRUEvictor[K, V]) Evict() (K, *types.CacheItem[K, V], bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element := l.list.Back(); element != nil {
		l.list.Remove(element)
		item := element.Value.(*types.CacheItem[K, V])
		delete(l.keys, item.Key)
		return item.Key, item, true
	}

	var zero K
	return zero, nil, false
}