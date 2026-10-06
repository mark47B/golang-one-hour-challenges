package main

import ( 
//"runtime"
// "time"
//  "fmt"
// "context"
// _ "log"
// "errors"
// _ "net/http/pprof"
// "net/http"
	"sync"
	// _ "sync/atomic"
)

type Item struct {
	prevItem *Item
	nextItem *Item
	value int
	key int
}

type LRUCache struct {
    head *Item
	tail *Item
	capacity int
	items map[int]*Item

	mu sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity < 1 {
		return nil
	}
	return &LRUCache{
		head: nil,
		tail: nil,
		capacity: capacity,
		items: make(map[int]*Item),
	}

}

func (c *LRUCache) Get(key int) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.get(key)

}

func (c *LRUCache) get(key int) (int, bool) {
	if item, ok := c.items[key]; !ok {
		return 0, false
	} else {
		if item == c.head {
			return item.value, true
		}

		if item == c.tail {
			c.tail = item.prevItem
		}

		prev := item.prevItem
		next := item.nextItem
		if next != nil {
			next.prevItem = prev
		}
		prev.nextItem = next
			
		item.nextItem = c.head
		item.prevItem = nil

		c.head.prevItem = item
		c.head = item
		return item.value, true
	}
}

func (c *LRUCache) Put(key int, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		item.value = value
		c.get(key)
		return
	}

	if len(c.items) >= c.capacity {
		delete(c.items, c.tail.key)
		if c.head == c.tail {
			c.head = nil
			c.tail = nil
		} else {
			prevTail := c.tail.prevItem
			prevTail.nextItem = nil
			c.tail = prevTail
		}
	}

	newItem := &Item{
		prevItem: nil,
		nextItem: c.head,
		value: value,
		key: key,
	}
	c.items[key] = newItem
	if c.head == nil {
		c.tail = newItem
	} else {
		c.head.prevItem = newItem
	}
	c.head = newItem
	return
}

func main(){
	_ = NewLRUCache(4)
}
