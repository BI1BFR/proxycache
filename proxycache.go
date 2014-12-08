//package proxycache is a key-value caching library.
package proxycache

import (
	"encoding/json"

	"github.com/huangml/proxycache/cache"
	"github.com/huangml/proxycache/proxy"
)

// ProxyCache is an in-memory key-value cache, and a database access proxy.
type ProxyCache struct {
	cache  *cache.Cache
	buffer *cache.Buffer
	saver  *proxy.Saver
	loader *proxy.Loader
}

// New creates a ProxyCache.
func New(p proxy.Proxy, maxEntry, saverProc, loaderProc int) *ProxyCache {
	c := cache.NewCache(maxEntry)
	b := cache.NewBuffer()
	s := proxy.NewSaver(p, saverProc, b)
	l := proxy.NewLoader(p, loaderProc)

	return &ProxyCache{
		cache:  c,
		buffer: b,
		saver:  s,
		loader: l,
	}
}

// Get retrieves data from ProxyCache.
// If provided key is not found in cache, data will be loaded by calling Proxy's
// Load method.
func (p *ProxyCache) Get(key string) []byte {
	entry := p.cache.Get(key)
	if entry != nil {
		return entry.Value
	}

	entry = p.buffer.Get(key)
	if entry != nil {
		return entry.Value
	}

	val, ok := p.loader.Load(key)
	if ok {
		p.cache.Put(&cache.Entry{key, val})
	}
	return val
}

// Put puts data into ProxyCache.
// Data will be saved asynchronously by calling Proxy's Save method.
func (p *ProxyCache) Put(key string, value []byte, ttw int64) {
	entry := &cache.Entry{key, value}
	p.cache.Put(entry)
	p.buffer.Put(entry, ttw)
}

// SetMaxEntry sets Cache's maxEntry.
func (p *ProxyCache) SetMaxEntry(maxEntry int) {
	p.cache.SetMaxEntry(maxEntry)
}

// SetLoadMaxProc sets Loader's maxProc.
func (p *ProxyCache) SetLoadMaxProc(maxProc int) {
	p.loader.SetMaxProc(maxProc)
}

// SetSaveProc sets the number of Saver's workers.
func (p *ProxyCache) SetSaveProc(proc int) {
	p.saver.SetMaxProc(proc)
}

// Status is used for runtime performance profiling.
type Status struct {
	cache.CacheStatus
	cache.BufferStatus
	proxy.LoaderStatus
	proxy.SaverStatus
}

// Status returns ProxyCache's runtime performance status.
func (p *ProxyCache) Status() []byte {
	s := Status{
		CacheStatus:  p.cache.Status(),
		BufferStatus: p.buffer.Status(),
		LoaderStatus: p.loader.Status(),
		SaverStatus:  p.saver.Status(),
	}

	b, _ := json.Marshal(s)
	return b
}
