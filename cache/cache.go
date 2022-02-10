package cache

import (
	"container/list"
	"fmt"
	"sync"
)

type Cache struct {
	locker sync.Mutex
	queue  *list.List
	cache  map[string]*list.Element
	size   int
}

type entry struct {
	Key   string
	Value interface{}
}

func NewCache(size int) *Cache {
	c := &Cache{
		queue: list.New(),
		cache: make(map[string]*list.Element),
		size:  size,
	}
	return c
}

func (c *Cache) Get(key string) interface{} {
	c.locker.Lock()
	defer c.locker.Unlock()
	element, ok := c.cache[key]
	if !ok || element == nil {
		return nil
	}
	en := element.Value.(*entry)
	data := en.Value

	// 访问之后，元素就移动到队列头部
	c.queue.MoveToFront(element)
	return data
}

func (c *Cache) Set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	element, ok := c.cache[key]
	if ok {
		en := element.Value.(*entry)
		en.Key = key
		en.Value = value
		return
	}

	// 不存在则建新entry
	en := &entry{Key: key, Value: value}
	newElem := c.queue.PushFront(en)
	c.cache[key] = newElem

	// 如果队列超过长度，则释放一部分
	if c.size < c.queue.Len() {
		tail := c.queue.Back()
		// 从队列中移除
		c.queue.Remove(tail)
		// 从map中移除
		remove := tail.Value.(*entry)
		fmt.Println("release", remove.Key)
		k := remove.Key
		delete(c.cache, k)
	}
}
