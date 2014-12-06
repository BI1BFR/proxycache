package cache

import (
	"sync"
	"time"

	"github.com/huangml/proxycache/priority-queue/pq"
)

type Buffer struct {
	entries map[string]*Entry
	q       *pq.PriorityQueue
	mtx     sync.Mutex

	ch   chan *Entry
	cond *sync.Cond
}

func NewBuffer() *Buffer {
	b := &Buffer{
		entries: make(map[string]*Entry),
		q:       pq.New(),
		ch:      make(chan *Entry),
	}
	b.cond = sync.NewCond(&b.mtx)

	go func() {
		for {
			if entry := func() *Entry {
				b.mtx.Lock()
				defer b.mtx.Unlock()

				for b.q.Len() == 0 {
					b.cond.Wait()
				}

				return b.entries[b.q.Pop().(string)]

			}(); entry != nil {
				b.ch <- entry
			}
		}
	}()

	return b
}

func (b *Buffer) SaveChan() <-chan *Entry {
	return b.ch
}

func (b *Buffer) Get(key string) *Entry {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	e, _ := b.entries[key]
	return e
}

func (b *Buffer) Put(entry *Entry, ttw int64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.entries[entry.Key] = entry
	priority := time.Now().Unix() + ttw

	if oldPriority, ok := b.q.Priority(entry.Key); ok && oldPriority <= priority {
		return
	}

	b.q.Push(entry.Key, priority)
	b.cond.Signal()
}

func (b *Buffer) OnSave(entry *Entry, ok bool) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if ok {
		if _, needSave := b.q.Priority(entry.Key); !needSave {
			delete(b.entries, entry.Key)
		}
	} else {
		if _, ok := b.entries[entry.Key]; !ok {
			b.entries[entry.Key] = entry
		}
		priority := time.Now().Unix() + 1
		b.q.Push(entry.Key, priority)
		b.cond.Signal()
	}
}

func (b *Buffer) Len() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return len(b.entries)
}
