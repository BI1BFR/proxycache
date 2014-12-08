// package cache defines containers for in-memory storage.
package cache

import (
	"sync"

	"github.com/huangml/proxycache/lru"
)

// Cache is an LRU cache.
// It auto removes least-recently-used entry when cache is full.
type Cache struct {
	maxEntry int
	entries  map[string]*Entry
	use      *lru.LRU
	mtx      sync.Mutex
}

// NewCache creates a new Cache.
// If maxEntry is 0, the cache has no limit size.
func NewCache(maxEntry int) *Cache {
	return &Cache{
		maxEntry: maxEntry,
		entries:  make(map[string]*Entry),
		use:      lru.New(),
	}
}

// SetMaxEntry setup a new maxEntry.
// Extra entries will be removed immediately.
func (c *Cache) SetMaxEntry(maxEntry int) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.maxEntry = maxEntry
	c.checkMaxEntryWithLock()
}

// Get looks up entry by a key.
// It marks the key as recently-used.
func (c *Cache) Get(key string) *Entry {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if e, ok := c.entries[key]; ok {
		c.use.Touch(key)
		return e
	} else {
		return nil
	}
}

// Put puts an entry to the cache.
// It marks the key as rencently-used.
func (c *Cache) Put(entry *Entry) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.entries[entry.Key] = entry
	c.use.Touch(entry.Key)
	c.checkMaxEntryWithLock()
}

// MaxEntry returns maxEntry of the cache.
func (c *Cache) MaxEntry() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.maxEntry
}

func (c *Cache) checkMaxEntryWithLock() {
	for c.maxEntry > 0 && len(c.entries) > c.maxEntry {
		if k, ok := c.use.Pop().(string); ok {
			delete(c.entries, k)
		} else {
			break
		}
	}
}

// CacheStatus is used for runtime performance profiling.
type CacheStatus struct {
	MaxEntry  int `json:"maxEntry"`
	CacheSize int `json:"cacheSize"`
}

// Status returns Cache's runtime performance status.
func (c *Cache) Status() CacheStatus {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return CacheStatus{
		MaxEntry:  c.maxEntry,
		CacheSize: len(c.entries),
	}
}
