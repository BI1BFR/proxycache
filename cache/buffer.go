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

	out    chan *Entry
	hasOut *sync.Cond
}

// NewBuffer creates a Buffer.
// It pumps entries to a channel order by their priorities.
func NewBuffer() *Buffer {
	b := &Buffer{
		entries: make(map[string]*Entry),
		q:       pq.New(),
		out:     make(chan *Entry),
	}
	b.hasOut = sync.NewCond(&b.mtx)

	go func() {
		for {
			if entry := func() *Entry {
				b.mtx.Lock()
				defer b.mtx.Unlock()

				for b.q.Len() == 0 {
					b.hasOut.Wait()
				}

				return b.entries[b.q.Pop().(string)]

			}(); entry != nil {
				b.out <- entry
			}
		}
	}()

	return b
}

// Entries returns a read only channel.
// All entries need to be save will be sent into this channel.
func (b *Buffer) Entries() <-chan *Entry {
	return b.out
}

// Get looks up an entry by the provided key.
func (b *Buffer) Get(key string) *Entry {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	e, _ := b.entries[key]
	return e
}

// Put puts an entry to Buffer.
// The entry will be pumped into out channel, orderd by priority.
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
	b.hasOut.Signal()
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
		b.hasOut.Signal()
	}
}

// BufferStatus is used for runtime performance profiling.
type BufferStatus struct {
	BufferSize int `json:"bufferSize"`
}

// Status returns Buffer's runtime performance status.
func (b *Buffer) Status() BufferStatus {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return BufferStatus{
		BufferSize: len(b.entries),
	}
}
