package proxy

import (
	"sync"

	"github.com/huangml/proxycache/cache"
)

type ProxySaver interface {
	Save(key string, value []byte) (ok bool)
}

type Buffer interface {
	SaveChan() <-chan *cache.Entry
	OnSave(entry *cache.Entry, ok bool)
}

type Saver struct {
	p      ProxySaver
	buffer Buffer
	*proc

	mtx      sync.Mutex // lock inFlight
	inFlight map[string]struct{}
}

func NewSaver(p ProxySaver, maxProc int, buffer Buffer) *Saver {
	s := &Saver{
		p:        p,
		buffer:   buffer,
		proc:     newProc(maxProc),
		inFlight: make(map[string]struct{}),
	}

	pipe := make(chan *cache.Entry)

	go func() {
		for {
			entry := <-s.buffer.SaveChan()
			s.mtx.Lock()
			if _, saving := s.inFlight[entry.Key]; saving {
				s.mtx.Unlock()
				s.buffer.OnSave(entry, false)
			} else {
				s.inFlight[entry.Key] = struct{}{}
				s.mtx.Unlock()
				pipe <- entry
			}

		}
	}()

	go func() {
		for {
			<-s.start
			go func() {
				for {
					select {
					case <-s.quit:
						return
					case entry := <-pipe:
						ok := s.p.Save(entry.Key, entry.Value)
						s.buffer.OnSave(entry, ok)
						s.mtx.Lock()
						delete(s.inFlight, entry.Key)
						s.mtx.Unlock()
					}
				}
			}()
		}
	}()

	return s
}
