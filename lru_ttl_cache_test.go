package gocache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRUTtlCache(t *testing.T) {
	fmt.Println("test lru ttl cache")
	gocache := NewGoCache(nil, nil, 3, 5, 0)
	assert.NotNil(t, gocache)

	value, err := gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err, "key golang not found")

	err = gocache.cache.Set("golang", "golang")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)

	time.Sleep(time.Second * 6)
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

	gocache = NewGoCache(valueGetter, valueMutiGetter, 3, 5, 0)
	assert.NotNil(t, gocache)

	value, err = gocache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)

	value, err = gocache.cache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)

	values, err := gocache.MGet([]string{"c1", "c2", "c3"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}

	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)

	err = gocache.cache.Set("c4", "c4")
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)
	err = gocache.cache.Set("c5", "c5")
	assert.Nil(t, err)
	time.Sleep(time.Second * 3)

	value, err = gocache.cache.Get("c4")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("c5")
	assert.Equal(t, value, "c5")
	assert.Nil(t, err)
}
