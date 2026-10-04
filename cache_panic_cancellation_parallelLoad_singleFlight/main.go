package main

import ( 
//"runtime"
// "time"
//  "fmt"
"context"
// "log"
// _ "net/http/pprof"
// "net/http"
	"sync"
)

type loader struct {
	value string
	err error
	panicVal any
	ready chan struct{}
}

type Cache struct {
    mu   sync.RWMutex
    data map[string]string
	
	loaders map[string]*loader
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
		loaders: make(map[string]*loader),
    }
}

func (c *Cache) runLoader(key string, l *loader, load func()(string, error)) {
	defer func(){
		if r := recover(); r != nil {
			l.panicVal = r
		}
		c.mu.Lock()
		delete(c.loaders, key)
		close(l.ready)
		c.mu.Unlock()
	}()
	newVal, err := load() // Для отмены через контекст load должен принимать ctx(пока такого нет - нет возможности отменить долгую операцию)
	l.value = newVal
	l.err = err
	if err != nil {
		return
	}
	c.mu.Lock()
	c.data[key] = newVal
	c.mu.Unlock()
}

func (c *Cache) Get(ctx context.Context, key string, load func() (string, error)) (string, error) {
	c.mu.RLock()
	val, ok := c.data[key]
	c.mu.RUnlock()
	if ok {
		return val, nil
	}

	c.mu.Lock()
	if val, ok := c.data[key]; ok {
		c.mu.Unlock()
		return val, nil
	}

	l, ok := c.loaders[key]
	if !ok {
		l = &loader{
			ready: make(chan struct{}),
		}
		c.loaders[key] = l
		go c.runLoader(key, l, load)
	}
	c.mu.Unlock()
	select {
	case <-l.ready:
		if l.panicVal != nil {
			panic(l.panicVal)
		}
		return l.value, l.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}


func main(){
	cache := NewCache()
	_ = cache
}
