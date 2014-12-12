//package pq implements a priority queue.
package pq

import "container/heap"

// PriorityQueue holds values. Values are auto orderd by their priorities.
type PriorityQueue struct {
	heap  pqHeap
	index map[interface{}]*pqEntry
}

// New creates a priority queue.
func New() *PriorityQueue {
	return &PriorityQueue{
		index: make(map[interface{}]*pqEntry),
	}
}

// Priority retrieves values's priority.
func (p *PriorityQueue) Priority(value interface{}) (priority int64, ok bool) {
	if e, ok := p.index[value]; ok {
		priority, ok = e.priority, true
	}
	return
}

// Push add a value with priority to queue.
func (p *PriorityQueue) Push(value interface{}, priority int64) {
	if e, ok := p.index[value]; ok {
		e.priority = priority
		heap.Fix(&p.heap, e.index)
	} else {
		e := &pqEntry{
			value:    value,
			priority: priority,
		}
		heap.Push(&p.heap, e)
		p.index[value] = e
	}
}

// Pop pops out value with minimal priority. It returns nil if no value exists.
func (p *PriorityQueue) Pop() interface{} {
	if e := heap.Pop(&p.heap).(*pqEntry); e != nil {
		delete(p.index, e.value)
		return e.value
	} else {
		return nil
	}
}

// Remove removes the provided value from queue.
func (p *PriorityQueue) Remove(value interface{}) {
	if e, ok := p.index[value]; ok {
		heap.Remove(&p.heap, e.index)
		delete(p.index, value)
	}
}

// Len returns queue size.
func (p *PriorityQueue) Len() int {
	return len(p.heap)
}
