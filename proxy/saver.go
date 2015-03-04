package proxy

import (
	"sync"
	"time"

	"github.com/huangml/proxycache/cache"
)

// ProxySaver is the interface wraps the Save method.
type ProxySaver interface {
	Save(key string, value []byte) (ok bool)
}

// Buffer is the interface of a save buffer.
type Buffer interface {
	Entries() <-chan *cache.Entry
	OnSave(entry *cache.Entry, ok bool)
}

// Saver calls ProxySaver concurrently to save entries from Buffer.
type Saver struct {
	*proc

	mtx      sync.Mutex
	inFlight map[string]struct{}
}

// NewSaver creates a Saver.
// It starts a group of workers do the saving.
// The number of workers is specified by parameter maxProc.
func NewSaver(p ProxySaver, maxProc int, buffer Buffer) *Saver {
	s := &Saver{
		proc:     newProc(maxProc),
		inFlight: make(map[string]struct{}),
	}

	pipe := make(chan *cache.Entry)

	// Fetch entries from Buffer and redirect to pipe.
	// If the given key is in saving, simply reports saving fail.
	go func() {
		for {
			entry := <-buffer.Entries()
			s.mtx.Lock()
			if _, saving := s.inFlight[entry.Key]; saving {
				s.mtx.Unlock()
				buffer.OnSave(entry, false)
			} else {
				s.inFlight[entry.Key] = struct{}{}
				s.mtx.Unlock()
				pipe <- entry
			}

		}
	}()

	// Fetch entries from pipe and do the real saving.
	go func() {
		for {
			// start a worker when get a start signal
			<-s.start
			go func() {
				for {
					select {
					// stop current worker when get a quit signal
					case <-s.quit:
						return
					case entry := <-pipe:
						ok := p.Save(entry.Key, entry.Value)
						buffer.OnSave(entry, ok)

						// remove from inFlight
						s.mtx.Lock()
						delete(s.inFlight, entry.Key)
						s.mtx.Unlock()

						if !ok {
							time.Sleep(time.Second)
						}
					}
				}
			}()
		}
	}()

	return s
}

// SaverStatus is used for runtime performance profiling.
type SaverStatus struct {
	MaxSaverProc int `json:"maxSaverProc"`
	SaverProc    int `json:"saverProc"`
	InflightSave int `json:"inflightSave"`
}

// Status returns Saver's runtime performance status.
func (s *Saver) Status() SaverStatus {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.proc.mtx.Lock()
	defer s.proc.mtx.Unlock()

	return SaverStatus{
		MaxSaverProc: s.proc.maxProc,
		SaverProc:    s.proc.maxProc - len(s.proc.start),
		InflightSave: len(s.inFlight),
	}
}
