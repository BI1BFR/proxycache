package volume

import (
	"sync"
	"time"

	"github.com/huangml/proxycache/priority-queue/pq"
)

type Buffer struct {
	entries map[string]*Entry
	q       *pq.PriorityQueue
	mtx     sync.Mutex

	saveChan chan *Entry
	saveCond *sync.Cond
}

func NewBuffer() *Buffer {
	b := &Buffer{
		entries:  make(map[string]*Entry),
		q:        pq.New(),
		saveChan: make(chan *Entry),
	}
	b.saveCond = sync.NewCond(&b.mtx)

	go func() {
		for {
			if entry := func() *Entry {
				b.mtx.Lock()
				defer b.mtx.Unlock()

				for b.q.Len() == 0 {
					b.saveCond.Wait()
				}

				return b.entries[b.q.Pop().(string)]

			}(); entry != nil {
				b.saveChan <- entry
			}
		}
	}()

	return b
}

func (b *Buffer) SaveChan() <-chan *Entry {
	return b.saveChan
}

func (b *Buffer) Get(key string) *Entry {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	e, _ := b.entries[key]
	return e
}

func (b *Buffer) Set(entry *Entry, ttw int64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.entries[entry.Key] = entry
	priority := time.Now().Unix() + ttw

	if oldPriority, ok := b.q.Priority(entry.Key); ok && oldPriority <= priority {
		return
	}

	b.q.Push(entry.Key, priority)
	b.saveCond.Signal()
}

func (b *Buffer) OnSave(key string, ok bool) {
	if ok {
		b.onSaveOK(key)
	} else {
		b.onSaveFail(key)
	}
}

func (b *Buffer) onSaveOK(key string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, needSave := b.q.Priority(key); !needSave {
		delete(b.entries, key)
	}
}

func (b *Buffer) onSaveFail(key string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	priority := time.Now().Unix()
	b.q.Push(key, priority)
	b.saveCond.Signal()
}

func (b *Buffer) Len() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return len(b.entries)
}
