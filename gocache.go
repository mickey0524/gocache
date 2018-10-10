package gocache

type Cache interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
}

type goCache struct {
	cache           Cache
	valueGetter     func(key string) (interface{}, error)
	valueMutiGetter func(keys []string) (map[string]interface{}, error)
}

func NewGoCache(valueGetter func(string) (interface{}, error),
	valueMutiGetter func([]string) (map[string]interface{}, error),
	cap int,
	ttl int) *goCache {
	var cache Cache

	switch {
	case cap > 0 && ttl > 0:
		cache = initLRUTtlCache(cap, ttl)
	case cap > 0:
		cache = initLRUCache(cap)
	default:
		cache = initKvCache()
	}

	return &goCache{
		valueGetter:     valueGetter,
		valueMutiGetter: valueMutiGetter,
		cache:           cache,
	}
}

func (c *goCache) Get(key string) (interface{}, error) {
	value, err := c.cache.Get(key)
	if err != nil && c.valueGetter != nil {
		value, err = c.valueGetter(key)
	}

	return value, err
}

func (c *goCache) MGet(keys []string) (map[string]interface{}, error) {
	var notFoundKeys []string
	res := make(map[string]interface{}, len(keys))

	for _, key := range keys {
		value, err := c.cache.Get(key)
		if err != nil {
			notFoundKeys = append(notFoundKeys, key)
		} else {
			res[key] = value
		}
	}

	if len(notFoundKeys) > 0 {
		if c.valueMutiGetter != nil {
			kvs, err := c.valueMutiGetter(notFoundKeys)
			if err == nil {
				for k, v := range kvs {
					res[k] = v
				}
			} else {
				for _, k := range notFoundKeys {
					res[k] = nil
				}
			}
		} else if c.valueGetter != nil {
			var missKeys []string
			for _, k := range notFoundKeys {
				v, err := c.valueGetter(k)
				if err == nil {
					res[k] = v
				} else {
					missKeys = append(missKeys, k)
				}
			}
			if len(missKeys) > 0 {
				for _, k := range missKeys {
					res[k] = nil
				}
			}
		} else {
			for _, k := range notFoundKeys {
				res[k] = nil
			}
		}
	}

	return res, nil
}
