package cache

import (
	"sync"
	"time"
	"github.com/lokeshllkumar/cache-store/v1/pkg/policies"
)

// main struct for in-memory cache store
//
// generics used to ensure type-safety of keys and values
type Cache[K comparable, V any] struct {
	items          sync.Map
	evictor        policies.EvictionPolicy[K, V]
	mu             sync.RWMutex
	capacity       int
	cleanupStop    chan struct{}
	cleanupWG      sync.WaitGroup
	cleanupEnabled bool
	metrics        *Metrics
	onEvictHook    func(K, V)
}

// config for initializing the cache
type Config[K comparable, V any] struct {
	Capacity        int
	EvictionPolicy  PolicyType
	CleanupInterval time.Duration
}

// defines the eviction policy type
type PolicyType string

const PolicyLRU PolicyType = "LRU"
const PolicyFIFO PolicyType = "FIFO"
const PolicyLFU PolicyType = "LFU"

// creates and initializes a new cache with generics
func NewCache[K comparable, V any](cfg Config[K, V]) *Cache[K, V] {
	if cfg.Capacity <= 0 {
		return nil
	}

	c := &Cache[K, V]{
		capacity:       cfg.Capacity,
		cleanupStop:    make(chan struct{}),
		cleanupEnabled: cfg.CleanupInterval > 0,
		metrics:        &Metrics{},
	}

	// initializes the specified eviction policy
	var evictor policies.EvictionPolicy[K, V]
	switch cfg.EvictionPolicy {
	case PolicyLRU:
		evictor = policies.NewLRUEvictor[K, V]()
	case PolicyFIFO:
		evictor = policies.NewFIFOEvictor[K, V]()
	case PolicyLFU:
		evictor = policies.NewLFUEvictor[K, V]()
	default:
		// fallback to LRU if none is specified
		evictor = policies.NewLRUEvictor[K, V]()
	}
	c.evictor = evictor

	if c.cleanupEnabled {
		c.startCleaner(cfg.CleanupInterval)
	}

	return c
}
