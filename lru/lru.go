package lru

import "container/list"

type LRU struct {
	l     *list.List
	index map[interface{}]*list.Element
}

func New() *LRU {
	return &LRU{
		l:     list.New(),
		index: make(map[interface{}]*list.Element),
	}
}

func (lru *LRU) Touch(key interface{}) {
	e, ok := lru.index[key]
	if ok {
		lru.l.MoveToBack(e)
	} else {
		lru.index[key] = lru.l.PushBack(key)
	}
}

func (lru *LRU) Len() int {
	return lru.l.Len()
}

func (lru *LRU) Pop() interface{} {
	if e := lru.l.Front(); e != nil {
		lru.Remove(e.Value)
		return e.Value
	} else {
		return nil
	}
}

func (lru *LRU) Remove(key interface{}) {
	if e, ok := lru.index[key]; ok {
		lru.l.Remove(e)
		delete(lru.index, key)
	}
}
