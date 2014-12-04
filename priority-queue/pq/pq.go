package pq

import "container/heap"

type PriorityQueue struct {
	heap  pqHeap
	index map[interface{}]*pqEntry
}

func New() *PriorityQueue {
	return &PriorityQueue{
		index: make(map[interface{}]*pqEntry),
	}
}

func (p *PriorityQueue) Priority(value interface{}) (priority int64, ok bool) {
	if e, ok := p.index[value]; ok {
		return e.priority, true
	} else {
		return 0, false
	}
}

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

func (p *PriorityQueue) Pop() interface{} {
	if e := heap.Pop(&p.heap).(*pqEntry); e != nil {
		return e.value
	} else {
		return nil
	}
}

func (p *PriorityQueue) Remove(value interface{}) {
	if e, ok := p.index[value]; ok {
		heap.Remove(&p.heap, e.index)
		delete(p.index, value)
	}
}

func (p *PriorityQueue) Len() int {
	return len(p.heap)
}
