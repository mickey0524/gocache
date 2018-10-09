package gocache

import (
	"container/list"
	"fmt"
	"sync"
)

type lruCacheIdx struct {
	value   interface{}
	pointer *list.Element
}

type lruCache struct {
	cap       int64
	index     map[string]*lruCacheIdx
	cacheList *list.List
	sync.RWMutex
}

func initLRUCache(cap int64) *lruCache {
	return &lruCache{
		cap:       cap,
		index:     make(map[string]*lruCacheIdx),
		cacheList: list.New(),
	}
}

func (c *lruCache) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	if idx, ok := c.index[key]; ok {
		return idx.value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *lruCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	if idx, ok := c.index[key]; ok {
		idx.value = value
		c.cacheList.MoveToBack(idx.pointer)
		return nil
	}

	if int64(c.cacheList.Len()) == c.cap {
		removeKey := c.cacheList.Front().Value.(string)
		c.remove(removeKey)
	}

	c.index[key] = &lruCacheIdx{
		value:   value,
		pointer: c.cacheList.PushBack(key),
	}

	return nil
}

func (c *lruCache) remove(key string) error {
	pointer := c.index[key].pointer
	delete(c.index, key)
	c.cacheList.Remove(pointer)

	return nil
}
