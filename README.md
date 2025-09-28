# cache-store

An in-memory cache client library for Go applications. Designed to be modular and concurrently safe, supporting multiple eviction policies and providing a clean, extensible API.

## Features

- Generics: Offers compile-time type safety for cache keys and values to prevent runtime errors
- Multiple Eviction Policies: Supports LRU (Least Recently Used), FIFO (First-In, First-Out), and LFU (Least Frequently Used) policies
- Concurrent Safety: Utilized `sync.Mutex` and `sync.Map` to ensure thread-safe access from multiple goroutines
- Auto Expiration: A dedicated background goroutine handle the asynchronous cleanup of expired items
- Metrics and Monitoring: Built-in support for tracking metrics like cache hits, misses, and evictions
- Event Hooks: Allows registration of callback functions to react to specific cache events, such as item eviction or expiration
- Graceful Shutdown: Supports the safe termination of background goroutines

## Getting Started

- Installation
    - Get the library using the `go get` command
    ```bash
    go get github.com/lokeshllkumar/cache-store/v1
    ```
- Configuration
    - Import the necessary packages and define your cache configuration using generics for type safety
    ```go
    package main

    import (
        "time"
        "github.com/lokeshllkumar/cache-store/v1/pkg/core"
        "github.com/lokeshllkumar/cache-store/v1/pkg/policies"
    )

    func main() {
        // 1. configure the cache: string keys and int values (for example)
        cfg := core.Config[string, int]{
            Capacity:        500,
            EvictionPolicy:  policies.PolicyLRU, // choose your policy
            CleanupInterval: 30 * time.Second,   // set to 0 to disable cleaner
        }

        // 2. initialize the Cache
        cacheStore := cache.NewCache(cfg)
        
        // 3. always defer Stop() to ensure the cleaner goroutine shuts down gracefully
        defer cacheStore.Stop() 
        
        // ... your application logic
    }
    ```

## Usage

- Setting and Retrieving Items
    - Use the `Set` method to add an item with a specific Time-To-Live (TTL), and `Get` to retrieve it
    ```go
        // Set an item that expires in 5 minutes
    cacheStore.Set("item1", 42, 5 * time.Minute)

    // Retrieve the item
    value, found := cacheStore.Get("item1")

    if found {
        fmt.Printf("item 1's value: %d\n", value)
    } else {
        fmt.Println("Cache Miss: Item not found or expired.")
    }
    ```

- Monitoring with Metrics
    - Use the `Metrics` method to get a snapshot of the cache's performance counters
    ```go
    metrics := cacheStore.Metrics()
    fmt.Printf("Cache Performance:\n")
    fmt.Printf("- Hits: %d\n", metrics.Hits)
    fmt.Printf("- Misses: %d", metrics.Misses)
    fmt.Printf("- Evictions: %d\n", metrics.Evictions)
    ```

- Handling Eviction Events
    - Register a hook function using `OnEvict` to execute custom logic whenever an item is removed from the cache due to eviction or expiration
    ```go
    // log the evicted key and value to an external system
    cacheStore.OnEvict(func(key string, value int) {
        fmt.Printf("Eviction Hook Fired: Key '%s' (Item value: %d) was removed.\n", key, value)
    })

    // ... set operations that trigger an eviction will call this hook
    ```