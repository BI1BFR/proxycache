package proxy

import "github.com/huangml/proxycache/volume"

type ProxySaver interface {
	Save(key string, value []byte) (ok bool)
}

type Saver struct {
	p       ProxySaver
	entries <-chan *volume.Entry
	*proc
}

func NewSaver(p ProxySaver, maxProc int, ch <-chan *volume.Entry) *Saver {
	s := &Saver{
		p:       p,
		entries: ch,
		proc:    newProc(maxProc),
	}

	go func() {
		for {
			<-s.start
			go s.save()
		}
	}()

	return s
}

func (s *Saver) save() {
	for {
		select {
		case <-s.quit:
			return
		case <-s.entries:
		}
	}
}
