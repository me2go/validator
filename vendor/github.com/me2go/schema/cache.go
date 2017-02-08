package schema

import "sync"

type Cache interface {
	Get(string) interface{}
	Put(string, interface{})
}

func NewCache() Cache {
	return &cache{}
}

type cache struct {
	sync.RWMutex
	holder map[string]interface{}
}

func (c *cache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.holder[key]
	if !ok {
		return nil
	}
	return v
}

func (c *cache) Put(key string, v interface{}) {
	c.Lock()
	defer c.Unlock()
	c.holder[key] = v
}
