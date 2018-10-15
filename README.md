## gocache

golang编写的cache缓存库，实现了基本的kv缓存，LRU缓存，带过期时间的LRU缓存，LRUK缓存，带过期时间的LRUK缓存

### 使用方法

```golang
govendor init
govendor fetch githu.com/mickey0524/gocache

import "githu.com/mickey0524/gocache"

gocache := NewGoCache(valueGetter, valueMutiGetter, cap, ttl, k)
```

### 参数列表

| 参数             | 说明     | 类型 |
| ---------------- |:--------:|:--------:|
| valueGetter      | 从外部数据源获取一个key的value | nil/func(string)interface{}|
| valueMutiGetter  | 从外部数据源获取多个key的value | nil/func([]string)map[string]interface{}|
| cap    | gocache最多缓存多少个key      | 0代表不设置上限|
| ttl    | gocache的key在多少秒后失效    | 0代表不设置过期时间|
| k      | lruk模式中当访问多少次将其放入cache| 0代表不设置k |

## API

* NewGoCache(valueGetter, valueMutiGetter, cap, ttl, k)

    该方法会返回如下结构的struct的实例指针

    ```golang
    type gocache struct {
        cache           Cache
	    valueGetter     func(key string) (interface{}, error)
	    valueMutiGetter func(keys []string) (map[string]interface{}, error)  
    }
    ```

    cache是gocache按照你的参数帮你生成的cache，valueGetter和valueMutiGetter就是你的参数

* gocache.Get(key string) interface{}, error

    这是gocache的API，当key存在于cache中时，直接返回，如果不在且设置了valueGetter，会自动调用valueGetter同时更新cache

* gocache.MGet(keys []string) map[string]interface{}, error

    这是gocache的API，keys是多个key组成的slice，当keys中有不存在于cache的时候，gocache会先检查是否设置了valueMutiGetter，再检查valueGetter，同样会自动调用同时更新cache

* gocache.cache.Set(key string, value string) error

    直接操作cache，将一个kv写入gocache

* gocache.cache.Get(key string) interface{}, error

    直接操作cache，从gocache中读取一个key
