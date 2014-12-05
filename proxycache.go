package proxycache

import "github.com/huangml/proxycache/volume"

type ProxyCache struct {
	volume *volume.Volume
	proxy  Proxy
}

func New(proxy Proxy, maxEntry int) *ProxyCache {
	p := &ProxyCache{
		volume: volume.NewVolume(maxEntry),
		proxy:  proxy,
	}

	go func() {
		for {
			entry := <-p.volume.SaveChan()
			ok := p.proxy.Save(entry.Key, entry.Value)
			p.volume.OnSave(entry.Key, ok)
		}
	}()

	return p
}

func (p *ProxyCache) Get(key string) []byte {
	if value := p.volume.Get(key); value != nil {
		return value
	}

	value, ok := p.proxy.Load(key)
	if ok {
		p.volume.OnLoad(key, value)
		return value
	}

	return nil
}

func (p *ProxyCache) Set(key string, value []byte, ttw int64) {
	p.volume.Set(key, value, ttw)
}
