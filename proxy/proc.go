package proxy

import "sync"

const MaxOfMaxProc = 64

type proc struct {
	start   chan struct{}
	quit    chan struct{}
	maxProc int
	mtx     sync.Mutex
}

func newProc(maxProc int) *proc {
	p := &proc{
		start: make(chan struct{}, MaxOfMaxProc),
		quit:  make(chan struct{}, MaxOfMaxProc),
	}
	p.SetMaxProc(maxProc)
	return p
}

func (p *proc) SetMaxProc(maxProc int) {
	if maxProc > MaxOfMaxProc {
		maxProc = MaxOfMaxProc
	} else if maxProc < 0 {
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
