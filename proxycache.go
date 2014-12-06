package proxycache

import (
	"github.com/huangml/proxycache/cache"
	"github.com/huangml/proxycache/proxy"
)

type ProxyCache struct {
	cache  *cache.Cache
	buffer *cache.Buffer
	saver  *proxy.Saver
	loader *proxy.Loader
}

func New(p proxy.Proxy, maxEntry int) *ProxyCache {
	c := cache.NewCache(maxEntry)
	b := cache.NewBuffer()
	s := proxy.NewSaver(p, 1, b)
	l := proxy.NewLoader(p, 1)

	return &ProxyCache{
		cache:  c,
		buffer: b,
		saver:  s,
		loader: l,
	}
}

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

func (p *ProxyCache) Put(key string, value []byte, ttw int64) {
	entry := &cache.Entry{key, value}
	p.cache.Put(entry)
	p.buffer.Put(entry, ttw)
}

func (p *ProxyCache) SetLoadMaxProc(maxProc int) {
	p.loader.SetMaxProc(maxProc)
}

func (p *ProxyCache) SetSaveMaxProc(maxProc int) {
	p.saver.SetMaxProc(maxProc)
}
