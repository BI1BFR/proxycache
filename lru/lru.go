// package lru implements an LRU queue.
package lru

import "container/list"

// LRU is an LRU queue.
type LRU struct {
	l     *list.List
	index map[interface{}]*list.Element
}

// New creates an LRU queue.
func New() *LRU {
	return &LRU{
		l:     list.New(),
		index: make(map[interface{}]*list.Element),
	}
}

// Touch marks a key as recently-used. The key will be created if not exists.
func (lru *LRU) Touch(key interface{}) {
	e, ok := lru.index[key]
	if ok {
		lru.l.MoveToBack(e)
	} else {
		lru.index[key] = lru.l.PushBack(key)
	}
}

// Len returns number of keys.
func (lru *LRU) Len() int {
	return lru.l.Len()
}

// Pop pops out the least-recently-used key. It returns nil if no key exists.
func (lru *LRU) Pop() interface{} {
	if e := lru.l.Front(); e != nil {
		lru.Remove(e.Value)
		return e.Value
	} else {
		return nil
	}
}

// Remove removes the provided key from LRU queue.
func (lru *LRU) Remove(key interface{}) {
	if e, ok := lru.index[key]; ok {
		lru.l.Remove(e)
		delete(lru.index, key)
	}
}
