package gocache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type lruTtlCacheIdx struct {
	value      interface{}
	expireTime int64
	visitPtr   *list.Element
	timePtr    *list.Element
}

type lruTtlCache struct {
	cap           int64
	ttl           int64
	index         map[string]*lruTtlCacheIdx
	cacheList     *list.List
	timeOrderList *list.List
	sync.Mutex
}

func initLRUTtlCache(cap int64, ttl int64) *lruTtlCache {
	return &lruTtlCache{
		cap:           cap,
		ttl:           ttl,
		index:         make(map[string]*lruTtlCacheIdx),
		cacheList:     list.New(),
		timeOrderList: list.New(),
	}
}

func (c *lruTtlCache) Get(key string) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.removeTimeOutKey()

	if idx, ok := c.index[key]; ok {
		c.cacheList.MoveToBack(idx.visitPtr)
		return idx.value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *lruTtlCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	c.removeTimeOutKey()

	if idx, ok := c.index[key]; ok {
		idx.value = value
		idx.expireTime = time.Now().Unix() + c.ttl
		c.cacheList.MoveToBack(idx.visitPtr)
		c.timeOrderList.MoveToBack(idx.timePtr)
		return nil
	}

	if int64(c.cacheList.Len()) == c.cap {
		frontEleKey := c.cacheList.Front().Value.(string)
		c.remove(frontEleKey)
		c.index[key] = &lruTtlCacheIdx{
			value:      value,
			expireTime: time.Now().Unix() + c.ttl,
			visitPtr:   c.cacheList.PushFront(key),
			timePtr:    c.timeOrderList.PushFront(key),
		}
	}

	return nil
}

func (c *lruTtlCache) remove(key string) error {
	idx := c.index[key]
	c.cacheList.Remove(idx.visitPtr)
	c.timeOrderList.Remove(idx.timePtr)
	delete(c.index, key)

	return nil
}

func (c *lruTtlCache) removeTimeOutKey() error {
	for {
		frontEleKey := c.timeOrderList.Front().Value.(string)
		frontIdx := c.index[frontEleKey]
		if frontIdx.expireTime > time.Now().Unix() {
			return nil
		}
		c.remove(frontEleKey)
	}
}
