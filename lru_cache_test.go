package gocache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	fmt.Println("test lru cache")
	gocache := NewGoCache(nil, nil, 3, 0, 0)
	assert.NotEqual(t, gocache, nil)

	value, err := gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err, "key golang not found")

	err = gocache.cache.Set("k1", "k1")
	assert.Nil(t, err)
	err = gocache.cache.Set("k2", "k2")
	assert.Nil(t, err)
	err = gocache.cache.Set("k3", "k3")
	assert.Nil(t, err)
	for _, k := range []string{"k1", "k2", "k3"} {
		value, err = gocache.cache.Get(k)
		assert.Equal(t, k, value)
	}

	err = gocache.cache.Set("k4", "k4")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("k1")
	assert.Nil(t, value)
	assert.Error(t, err)
	for _, k := range []string{"k2", "k3", "k4"} {
		value, err = gocache.cache.Get(k)
		assert.Equal(t, k, value)
	}

	valueGetter := func(key string) (interface{}, error) {
		return key, nil
	}

	valueMutiGetter := func(keys []string) (map[string]interface{}, error) {
		res := make(map[string]interface{}, len(keys))

		for _, k := range keys {
			res[k] = k
		}

		return res, nil
	}

	gocache = NewGoCache(valueGetter, valueMutiGetter, 3, 0, 0)
	value, err = gocache.Get("k1")
	assert.Equal(t, value, "k1")
	assert.Nil(t, err)

	values, err := gocache.MGet([]string{"k1", "k2", "k3"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}

	values, err = gocache.MGet([]string{"k4", "k5", "k6"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}
	value, err = gocache.cache.Get("k1")
	assert.Nil(t, value)
	value, err = gocache.cache.Get("k2")
	assert.Nil(t, value)
	value, err = gocache.cache.Get("k3")
	assert.Nil(t, value)
	value, err = gocache.cache.Get("k4")
	assert.Equal(t, value, "k4")
	value, err = gocache.cache.Get("k5")
	assert.Equal(t, value, "k5")
	value, err = gocache.cache.Get("k6")
	assert.Equal(t, value, "k6")
}
