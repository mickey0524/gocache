package gocache

import (
	"container/heap"
	"fmt"
	"sync"
)

type lfuCacheIdx struct {
	key   string
	value interface{}
	rev   int
	freq  int
	index int
}

type lfuHeap []*lfuCacheIdx

func (h lfuHeap) Len() int {
	return len(h)
}

func (h lfuHeap) Less(i, j int) bool {
	if h[i].freq == h[j].freq {
		return h[i].rev < h[j].rev
	}

	return h[i].freq < h[j].freq
}

func (h lfuHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *lfuHeap) Push(x interface{}) {
	length := len(*h)
	item := x.(*lfuCacheIdx)
	item.index = length
	*h = append(*h, item)
}

func (h *lfuHeap) Pop() interface{} {
	length := len(*h)
	heap := *h
	item := heap[length-1]
	*h = heap[0 : length-1]

	return item
}

func (h *lfuHeap) update(item *lfuCacheIdx) {
	heap.Fix(h, item.index)
}

type lfuCache struct {
	cap   int
	index map[string]*lfuCacheIdx
	rev   int
	heap  lfuHeap
	sync.Mutex
}

func initLFUCache(cap int) *lfuCache {
	if cap <= 0 {
		panic("cap can not less then 0")
	}
	return &lfuCache{
		cap:   cap,
		index: make(map[string]*lfuCacheIdx),
	}
}

func (c *lfuCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if idx, ok := c.index[key]; ok {
		c.rev++
		idx.freq++
		idx.rev = c.rev
		c.heap.update(idx)

		return idx.value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *lfuCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	c.rev++
	if idx, ok := c.index[key]; ok {
		idx.value = value
		idx.rev = c.rev
		idx.freq++
		c.heap.update(idx)

		return nil
	}

	if len(c.heap) == c.cap {
		removeItem := heap.Pop(&c.heap).(lfuCacheIdx)
		delete(c.index, removeItem.key)
	}

	newItem := &lfuCacheIdx{
		key:   key,
		value: value,
		rev:   c.rev,
		freq:  1,
	}
	heap.Push(&c.heap, newItem)
	c.index[key] = newItem

	return nil
}
