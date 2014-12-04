package pq

type pqEntry struct {
	value    interface{}
	priority int64
	index    int
}

type pqHeap []*pqEntry

func (h pqHeap) Len() int {
	return len(h)
}

func (h pqHeap) Less(i, j int) bool {
	return h[i].priority < h[j].priority
}

func (h pqHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *pqHeap) Push(x interface{}) {
	n := len(*h)
	entry := x.(*pqEntry)
	entry.index = n
	*h = append(*h, entry)
}

func (h *pqHeap) Pop() interface{} {
	old := *h
	n := len(old)
	entry := old[n-1]
	entry.index = -1
	*h = old[0 : n-1]
	return entry
}
