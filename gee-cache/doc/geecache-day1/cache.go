package geecacheday1

import (
	"geecache-day1/lru" // 在同一个包里
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {  // 这里的ByteView 是不是引入了外部包里的
	c.mu.Lock()
	defer c.mu.Unlock()  // 在函数返回之前先执行这个操作 在此之前 函数中其他代码都无法访问c.mu
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}