package main

import (
	"fmt"
	"sync"
)

type InMemoryDB struct {
	mtx  sync.RWMutex
	data map[string][]byte
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string][]byte),
	}
}

func (d *InMemoryDB) Load(key string) (value []byte, ok bool) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	fmt.Println("Loading: ", key)

	v, ok := d.data[key]
	return v, ok
}

func (d *InMemoryDB) Save(key string, value []byte) bool {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	fmt.Println("Saving: ", key)

	d.data[key] = value
	return true
}
