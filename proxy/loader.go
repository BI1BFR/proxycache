package proxy

import "sync"

// ProxyLoader is the interface wraps the Load method.
type ProxyLoader interface {
	Load(key string) (value []byte, ok bool)
}

// Loader provides method to load data by Proxy concurrently.
type Loader struct {
	p ProxyLoader
	*proc

	mtx      sync.Mutex
	inFlight map[string]*loadResult
}

// NewLoader creates a Loader.
// Parameter maxProc specifies the maximum number of goroutines call Load(),
// the excess will be blocked.
func NewLoader(p ProxyLoader, maxProc int) *Loader {
	l := &Loader{
		p:        p,
		proc:     newProc(maxProc),
		inFlight: make(map[string]*loadResult),
	}

	go func() {
		for {
			<-l.quit
			<-l.start
		}
	}()

	return l
}

type loadResult struct {
	done  chan struct{}
	value []byte
	ok    bool
}

// Load loads data by the provided key concurrently.
// Duplicate keys will be loaded only once.
func (l *Loader) Load(key string) ([]byte, bool) {
	l.mtx.Lock()
	if f, ok := l.inFlight[key]; ok {
		l.mtx.Unlock()
		<-f.done
		return f.value, f.ok
	}

	f := &loadResult{done: make(chan struct{})}
	l.inFlight[key] = f

	l.mtx.Unlock()

	<-l.start
	f.value, f.ok = l.p.Load(key)
	l.start <- struct{}{}
	close(f.done)

	l.mtx.Lock()
	delete(l.inFlight, key)
	l.mtx.Unlock()

	return f.value, f.ok
}

// ProcNum returns number of loading goroutines.
func (l *Loader) ProcNum() int {
	l.proc.mtx.Lock()
	defer l.proc.mtx.Unlock()

	return l.maxProc - len(l.start)
}
