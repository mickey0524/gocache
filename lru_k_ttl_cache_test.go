package gocache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRUKTtlCache(t *testing.T) {
	fmt.Println("test lru k ttl cache")
	gocache := NewGoCache(nil, nil, 3, 3, 2)
	assert.NotNil(t, gocache)

	gocache.cache.Set("golang", "golang")
	value, err := gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Equal(t, value, "golang")
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)

	gocache.cache.Set("golang", "golang")
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)
	time.Sleep(time.Second * 3)
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("golang")
	assert.Nil(t, value)
	assert.Error(t, err)

	for i := 0; i < 3; i++ {
		for _, k := range []string{"c1", "c2", "c3"} {
			if i == 0 {
				err = gocache.cache.Set(k, k)
				assert.Nil(t, err)
			}
			value, err = gocache.cache.Get(k)
			if i < 2 {
				assert.Nil(t, value)
				assert.Error(t, err)
			} else {
				assert.Equal(t, k, value)
				assert.Nil(t, err)
			}
		}
	}

	err = gocache.cache.Set("c4", "c4")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("c4")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("c4")
	assert.Nil(t, value)
	assert.Error(t, err)
	value, err = gocache.cache.Get("c4")
	assert.Equal(t, value, "c4")
	assert.Nil(t, err)
	value, err = gocache.cache.Get("c1")
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

	gocache = NewGoCache(valueGetter, valueMutiGetter, 3, 3, 2)

	values, err := gocache.MGet([]string{"c1", "c2", "c3"})
	assert.Nil(t, err)
	for k, v := range values {
		assert.Equal(t, k, v)
	}

	for i := 0; i < 3; i++ {
		for _, k := range []string{"c1", "c2", "c3"} {
			value, err = gocache.cache.Get(k)
			if i < 2 {
				assert.Nil(t, value)
				assert.Error(t, err)
			} else {
				assert.Equal(t, k, value)
				assert.Nil(t, err)
			}
		}
	}
}
