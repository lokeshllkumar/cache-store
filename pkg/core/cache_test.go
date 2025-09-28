package cache

import (
	"testing"
	"time"
)

// verifies that the cache can be created and handles invalid config
func TestNewCache(t *testing.T) {
	// test valid config
	cfg := Config[string, int]{
		Capacity: 10,
		EvictionPolicy: PolicyLRU,
		CleanupInterval: 5 * time.Minute,
	}
	c := NewCache(cfg)
	if c == nil {
		t.Fatal("NewCache retuned nil for a valid configuration")
	}

	// testing invalid capacity
	invalidCfg := Config[string, int]{
		Capacity: 0,
	}
	c = NewCache(invalidCfg)
	if c != nil {
		t.Fatal("NewCache did not return nil for invalid capacity")
	}
}

// validates basic set + get functionality
func TestSetGet(t *testing.T) {
	cfg := Config[string, int]{
		Capacity: 10,
		EvictionPolicy: PolicyLRU,
		CleanupInterval: 50 * time.Millisecond,
	}
	c := NewCache(cfg)
	defer c.Stop()

	// short-lived item initialized/set
	c.Set("short_lived_key", 1234, 10 * time.Millisecond)

	// waiting for time to expire
	time.Sleep(20 * time.Millisecond)
	_, found := c.Get("short_lived_key")
	if found {
		t.Error("Item did not expire after its TTL")
	}

	// waiting for the cleaner to run and remove the item
	time.Sleep(100 * time.Millisecond)
	_, found = c.Get("short_lived_key")
	if found {
		t.Error("Cleaner failed to remove expired item")
	}
}

// verifies LRU policy functionality
func TestLRUEviction(t *testing.T) {
	cfg := Config[string, int]{
		Capacity: 2,
		EvictionPolicy: PolicyLRU,
		CleanupInterval: 0, // disabling the cleaner
	}

	c := NewCache(cfg)

	// set 2 elements
	c.Set("a", 1, 0)
	c.Set("b", 2, 0)

	// get 1 element
	c.Get("a")

	// set new element to evict unaccessed element
	c.Set("c", 3, 0)

	// attempt to access element that should've been evicted
	_, found := c.Get("b")
	if found {
		t.Error("Expected key 'b' to be evicted by LRU policy")
	}

	// try to access elements that are in the cache
	_, found = c.Get("a")
	if !found {
		t.Error("Expected key 'a' to still be present in the cache")
	}
	_, found = c.Get("c")
	if !found {
		t.Error("Expected key 'c' to still be present in the cache")
	}
}

func TestAccessAcrossPolicies(t *testing.T) {
	policies := []PolicyType{PolicyLRU, PolicyFIFO, PolicyLFU}
	
	for _, policy := range policies {
		t.Run(string(policy), func(t *testing.T) {
			cfg := Config[string, string]{
				Capacity:       1,
				EvictionPolicy: policy,
			}
			c := NewCache(cfg)
			if c == nil {
				t.Fatalf("Failed to create cache with policy: %s", policy)
			}
			
			// set a new item
			c.Set("test_key", "test_value", 10 * time.Second)
			
			// check if the item is accessible
			val, found := c.Get("test_key")
			if !found {
				t.Errorf("Item not found for policy %s", policy)
			}
			
			if val != "test_value" {
				t.Errorf("Incorrect value retrieved for policy %s: expected %s, got %s", policy, "test_value", val)
			}
		})
	}
}