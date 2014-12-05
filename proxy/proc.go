package proxy

import "sync"

type proc struct {
	start   chan struct{}
	quit    chan struct{}
	maxProc int
	mtx     sync.Mutex
}

func newProc() *proc {
	return &proc{
		start: make(chan struct{}),
		quit:  make(chan struct{}),
	}
}

func (p *proc) SetMaxProc(maxProc int) {
	if maxProc < 0 {
		maxProc = 0
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	delta := maxProc - p.maxProc
	p.maxProc = maxProc

	go func() {
		if delta > 0 {
			for i := 0; i < delta; i++ {
				p.start <- struct{}{}
			}
		} else if delta < 0 {
			for i := 0; i < -delta; i++ {
				p.quit <- struct{}{}
			}
		}
	}()
}
