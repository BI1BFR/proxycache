package volume

import (
	"sync"

	"github.com/huangml/proxycache/lru"
)

type Cache struct {
	maxEntry int
	entries  map[string]*Entry
	use      *lru.LRU
	mtx      sync.Mutex
}

func NewCache(maxEntry int) *Cache {
	return &Cache{
		maxEntry: maxEntry,
		entries:  make(map[string]*Entry),
		use:      lru.New(),
	}
}

func (c *Cache) SetMaxEntry(maxEntry int) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.maxEntry = maxEntry
	c.checkMaxEntryWithLock()
}

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

func (c *Cache) Set(entry *Entry) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.entries[entry.Key] = entry
	c.use.Touch(entry.Key)
	c.checkMaxEntryWithLock()
}

func (c *Cache) OnLoad(key string, value []byte) {
	c.Set(&Entry{key, value})
}

func (c *Cache) Len() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return len(c.entries)
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
