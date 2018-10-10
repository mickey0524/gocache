package gocache

import (
	"container/list"
	"fmt"
	"sync"
)

type lruKCacheIdx struct {
	value   interface{}
	pointer *list.Element
}

type lruKCacheHistoryIdx struct {
	value     interface{}
	pointer   *list.Element
	touchTime int
}

type lruKCache struct {
	cacheIndex   map[string]*lruKCacheIdx
	cacheList    *list.List
	historyIndex map[string]*lruKCacheHistoryIdx
	historyList  *list.List
	cap          int
	k            int
	sync.Mutex
}

func initLRUKCache(cap int, k int) *lruKCache {
	return &lruKCache{
		cacheIndex:   make(map[string]*lruKCacheIdx),
		historyIndex: make(map[string]*lruKCacheHistoryIdx),
		cap:          cap,
		k:            k,
	}
}

func (c *lruKCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	if cacheIdx, ok := c.cacheIndex[key]; ok {
		c.cacheList.MoveToBack(cacheIdx.pointer)
		return cacheIdx.value, nil
	}

	if historyIdx, ok := c.historyIndex[key]; ok {
		historyIdx.touchTime++
		if historyIdx.touchTime == c.k {
			if c.cacheList.Len() == c.cap {
				c.removeCacheKey(c.cacheList.Front().Value.(string))
			}
			c.cacheIndex[key] = &lruKCacheIdx{
				value:   historyIdx.value,
				pointer: c.cacheList.PushBack(key),
			}
			c.removeHistoryKey(key)
		}
		c.historyList.MoveToBack(historyIdx.pointer)
		return historyIdx.value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *lruKCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	if cacheIdx, ok := c.cacheIndex[key]; ok {
		cacheIdx.value = value
		c.cacheList.MoveToBack(cacheIdx.pointer)
		return nil
	}

	if historyIdx, ok := c.historyIndex[key]; ok {
		historyIdx.value = value
		c.historyList.MoveToBack(historyIdx.pointer)
		return nil
	}

	if c.historyList.Len() == c.cap {
		c.removeHistoryKey(c.historyList.Front().Value.(string))
	}
	c.historyIndex[key] = &lruKCacheHistoryIdx{
		value:   value,
		pointer: c.historyList.PushBack(key),
	}

	return nil
}

func (c *lruKCache) removeHistoryKey(key string) error {
	historyIdx := c.historyIndex[key]
	c.historyList.Remove(historyIdx.pointer)
	delete(c.historyIndex, key)

	return nil
}

func (c *lruKCache) removeCacheKey(key string) error {
	cacheIdx := c.cacheIndex[key]
	c.cacheList.Remove(cacheIdx.pointer)
	delete(c.cacheIndex, key)

	return nil
}
