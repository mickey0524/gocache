package gocache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVCache(t *testing.T) {
	fmt.Println("test kv cache")
	gocache := NewGoCache(nil, nil, 0, 0, 0)
	assert.NotEqual(t, gocache, nil)

	value, err := gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err, "key golang not found")

	err = gocache.cache.Set("golang", "golang")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)

	values, err := gocache.MGet([]string{"golang", "cache"})
	assert.Nil(t, err)
	for k, v := range values {
		if k == "golang" {
			assert.Equal(t, k, v)
		} else {
			assert.Nil(t, v)
		}
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

	gocache = NewGoCache(valueGetter, valueMutiGetter, 0, 0, 0)

	value, err = gocache.Get("cache")
	assert.Equal(t, value, "cache")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("cache")
	assert.Equal(t, value, "cache")
	assert.Nil(t, err)

	values, err = gocache.MGet([]string{"c1", "c2"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}
	value, err = gocache.cache.Get("c1")
	assert.Equal(t, "c1", "c1")
	assert.Nil(t, err)
}
