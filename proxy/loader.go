package proxy

type ProxyLoader interface {
	Load(key string) (value []byte, ok bool)
}

type Loader struct {
	p ProxyLoader
	*proc
}

func NewLoader(p ProxyLoader, maxProc int) *Loader {
	l := &Loader{
		p:    p,
		proc: newProc(maxProc),
	}

	go func() {
		for {
			<-l.quit
			<-l.start
		}
	}()

	return l
}

func (l *Loader) Load(key string) ([]byte, bool) {
	<-l.start
	defer func() { l.start <- struct{}{} }()

	return l.p.Load(key)
}
