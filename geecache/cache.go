package geecache

import (
	"geecache/store"
	"sync"
)

const (
	LRU  = "lru"
	LFU  = "lfu"
	FIFO = "fifo"
)

type cache struct {
	cacheBytes int64
	strategy   string
	store      store.Store
	mu         sync.Mutex // todo: use Mutex will affect performance?
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store == nil {
		c.lazyInitStore(c.strategy, c.cacheBytes, nil)
	}
	c.store.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store == nil {
		return
	}
	if v, ok := c.store.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func (c *cache) lazyInitStore(strategy string, cacheBytes int64, onEvicted store.EvictedFunc) {
	var s store.Store
	switch strategy {
	case LRU:
		s = store.NewLRUCache(cacheBytes, onEvicted)
	case LFU:
		panic("LFU: Not implemented yet")
	case FIFO:
		panic("FIFO: Not implemented yet")
	default:
		panic("cache strategy not supported")
	}
	c.store = s
}
