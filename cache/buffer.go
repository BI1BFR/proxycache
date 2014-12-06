package cache

import (
	"sync"
	"time"

	"github.com/huangml/proxycache/priority-queue/pq"
)

// Buffer is a container which holds entries need to be save.
type Buffer struct {
	entries map[string]*Entry
	q       *pq.PriorityQueue
	mtx     sync.Mutex

	ch   chan *Entry
	cond *sync.Cond
}

// NewBuffer creates a Buffer.
// It pumps entries to a channel order by their priorities.
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

// SaveChan returns a read only channel.
// All entries need to be save will be filled into this channel.
func (b *Buffer) SaveChan() <-chan *Entry {
	return b.ch
}

// Get looks up an entry by the provided key.
func (b *Buffer) Get(key string) *Entry {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	e, _ := b.entries[key]
	return e
}

// Put puts an entry to Buffer.
// The entry will be pumped into SaveChan, orderd by priority.
// Entry's priority is `current unix epoch time + ttw`.
func (b *Buffer) Put(entry *Entry, ttw int64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.entries[entry.Key] = entry
	priority := time.Now().Unix() + ttw

	// if the key is already in queue, use the higher priority.
	if oldPriority, ok := b.q.Priority(entry.Key); ok && oldPriority <= priority {
		return
	}

	b.q.Push(entry.Key, priority)
	b.cond.Signal()
}

// OnSave handles entry saving result.
// If succesed, the entry will removed from Buffer.
// If failed, the entry will be pushed back to Buffer and saved later again.
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

// Len returns number of entries.
func (b *Buffer) Len() int {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return len(b.entries)
}
