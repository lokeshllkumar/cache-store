package cache

import (
	"sync/atomic"
)

// holds the performance data for the cache
type Metrics struct {
	Hits uint64
	Misses uint64
	Evictions uint64
	Items uint64
}

// increments hits counter (safer implementation)
func (m *Metrics) IncrementHits() {
	atomic.AddUint64(&m.Hits, 1)
}

// increments misses counter (safer implementation)
func (m *Metrics) IncrementMisses() {
	atomic.AddUint64(&m.Misses, 1)
}

// increments evictions counter (safer implementation)
func (m *Metrics) IncrementEvictions() {
	atomic.AddUint64(&m.Evictions, 1)
}