package lru

import (
	"container/list"
)

// Value 值接口，Len返回值多占的字节数
type Value interface {
	Len() int
}

type Cache struct {
	maxBytes int64 // 0表示不限制
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// 一条记录
type entry struct {
	key   string
	value Value
}

// New 实例化一个缓存
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找功能
func (c *Cache) Get(key string) (value Value, exists bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除功能，淘汰最近最少使用的key
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 更新内存占用
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 修改/更新
func (c *Cache) Add(key string, value Value) {
	// 如果key已经存在，则更新值， 并更新链表
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新内存占用
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 判断内存是否超限
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 获取当前缓存中的键值对数
func (c *Cache) Len() int {
	return c.ll.Len()
}
