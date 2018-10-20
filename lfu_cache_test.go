package gocache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFUCache(t *testing.T) {
	fmt.Println("test lfu cache")
	gocache := NewGoLFUCache(nil, nil, 3)
	assert.NotNil(t, gocache)

	err := gocache.cache.Set("golang", "golang")
	assert.Nil(t, err)
	value, err := gocache.cache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)
	for _, k := range []string{"c1", "c2", "c3"} {
		err = gocache.cache.Set(k, k)
		assert.Nil(t, err)
	}

	value, err = gocache.cache.Get("c1")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("c2")
	assert.Equal(t, value, "c2")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("c3")
	assert.Equal(t, value, "c3")
	assert.Nil(t, err)
	err = gocache.cache.Set("c1", "c1")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)

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

	gocache = NewGoLFUCache(valueGetter, valueMutiGetter, 3)
	assert.NotNil(t, gocache)

	value, err = gocache.Get("c1")
	assert.Equal(t, value, "c1")
	assert.Nil(t, err)

	values, err := gocache.MGet([]string{"k1", "k2", "k3"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}

	value, err = gocache.cache.Get("c1")
	assert.Nil(t, value)
	assert.Error(t, err)
}
