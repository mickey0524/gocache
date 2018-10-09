package gocache

import (
	"fmt"
	"sync"
)

type kvCache struct {
	data map[string]interface{}
	sync.RWMutex
}

func initKvCache() *kvCache {
	return &kvCache{
		data: make(map[string]interface{}),
	}
}

func (c *kvCache) Get(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	if value, ok := c.data[key]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("key %s not found", key)
}

func (c *kvCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	c.data[key] = value

	return nil
}
