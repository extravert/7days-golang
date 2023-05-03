package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element  // 键是字符串，值是双向链表中对应节点的指针。 这里就是entry
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes  接口算是方法的集成 
type Value interface {
	Len() int
}
// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll: 	   list.New(),
		cache: 	   make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 方法名前面括号表示是结构体Cache的方法 （而Go只是把C语言中的第一个参数放到方法前面而已，所以它并不是用来类型转换的，而是一个接收者，说明该方法属于哪个结构体。）
func (c *Cache) Get(key string) (value Value, ok bool) {
	// ele, ok := c.cache[key]
	// if ok {}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)  // 断言
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()  // back()函数 取到队首元素
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())  // 删除占用的空间 key和value长度
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if  ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())  // 这里还没看懂
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	} 
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

