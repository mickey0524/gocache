package gocache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type lruKTtlCacheIdx struct {
	value      interface{}
	expireTime int
	visitPtr   *list.Element
	timePtr    *list.Element
}

type lruKTtlCacheHistoryIdx struct {
	value      interface{}
	expireTime int
	visitPtr   *list.Element
	timePtr    *list.Element
	touchTime  int
}

type lruKTtlCache struct {
	cacheIndex           map[string]*lruKTtlCacheIdx
	cacheList            *list.List
	timeOrderCacheList   *list.List
	historyIndex         map[string]*lruKTtlCacheHistoryIdx
	historyList          *list.List
	timeOrderHistoryList *list.List
	cap                  int
	k                    int
	ttl                  int
	sync.Mutex
}

func initLRUKTtlCache(cap int, k int, ttl int) *lruKTtlCache {
	return &lruKTtlCache{
		cacheIndex:   make(map[string]*lruKTtlCacheIdx),
		historyIndex: make(map[string]*lruKTtlCacheHistoryIdx),
		cap:          cap,
		k:            k,
		ttl:          ttl,
	}
}

func (c *lruKTtlCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.removeTimeOutCacheKey()

	if cacheIdx, ok := c.cacheIndex[key]; ok {
		c.cacheList.MoveToBack(cacheIdx.visitPtr)
		return cacheIdx.value, nil
	}

	c.removeTimeOutHistoryKey()

	if historyIdx, ok := c.historyIndex[key]; ok {
		historyIdx.touchTime++
		if historyIdx.touchTime == c.k {
			if c.cacheList.Len() == c.cap {
				c.removeCacheKey(c.cacheList.Front().Value.(string))
			}
			c.cacheIndex[key] = &lruKTtlCacheIdx{
				value:      historyIdx.value,
				expireTime: int(time.Now().Unix()) + c.ttl,
				visitPtr:   c.cacheList.PushBack(key),
				timePtr:    c.timeOrderCacheList.PushBack(key),
			}
			c.removeHistoryKey(key)
		}
		c.historyList.MoveToBack(historyIdx.visitPtr)
		return historyIdx.value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *lruKTtlCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	c.removeTimeOutCacheKey()

	if cacheIdx, ok := c.cacheIndex[key]; ok {
		cacheIdx.value = value
		cacheIdx.expireTime = int(time.Now().Unix()) + c.ttl
		c.cacheList.MoveToBack(cacheIdx.visitPtr)
		c.timeOrderCacheList.MoveToBack(cacheIdx.timePtr)
		return nil
	}

	c.removeTimeOutHistoryKey()

	if historyIdx, ok := c.historyIndex[key]; ok {
		historyIdx.value = value
		historyIdx.expireTime = int(time.Now().Unix()) + c.ttl
		c.historyList.MoveToBack(historyIdx.visitPtr)
		c.timeOrderHistoryList.MoveToBack(historyIdx.timePtr)
		return nil
	}

	if c.historyList.Len() == c.cap {
		c.removeHistoryKey(c.historyList.Front().Value.(string))
	}
	c.historyIndex[key] = &lruKTtlCacheHistoryIdx{
		value:      value,
		expireTime: int(time.Now().Unix()) + c.ttl,
		visitPtr:   c.historyList.PushBack(key),
		timePtr:    c.timeOrderHistoryList.PushBack(key),
	}

	return nil
}

func (c *lruKTtlCache) removeHistoryKey(key string) error {
	historyIdx := c.historyIndex[key]
	c.historyList.Remove(historyIdx.visitPtr)
	c.historyList.Remove(historyIdx.timePtr)
	delete(c.historyIndex, key)

	return nil
}

func (c *lruKTtlCache) removeCacheKey(key string) error {
	cacheIdx := c.cacheIndex[key]
	c.cacheList.Remove(cacheIdx.visitPtr)
	c.cacheList.Remove(cacheIdx.timePtr)
	delete(c.cacheIndex, key)

	return nil
}

func (c *lruKTtlCache) removeTimeOutHistoryKey() error {
	length := c.timeOrderHistoryList.Len()

	for i := 0; i < length; i++ {
		frontKey := c.timeOrderHistoryList.Front().Value.(string)
		frontIdx := c.historyIndex[frontKey]
		if frontIdx.expireTime > int(time.Now().Unix()) {
			break
		}
		c.removeHistoryKey(frontKey)
	}

	return nil
}

func (c *lruKTtlCache) removeTimeOutCacheKey() error {
	length := c.timeOrderCacheList.Len()

	for i := 0; i < length; i++ {
		frontKey := c.timeOrderCacheList.Front().Value.(string)
		frontIdx := c.cacheIndex[frontKey]
		if frontIdx.expireTime > int(time.Now().Unix()) {
			break
		}
		c.removeCacheKey(frontKey)
	}

	return nil
}
