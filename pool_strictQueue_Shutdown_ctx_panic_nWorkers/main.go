package main

import ( 
//"runtime"
"time"
//  "fmt"
"context"
"log"
"errors"
// _ "net/http/pprof"
// "net/http"
	"sync"
)


type Pool struct {
    q chan func()error
	nWorkers int
	wg sync.WaitGroup
	active chan struct{}
	once sync.Once

	mu sync.RWMutex
	workersDone chan struct{}

}

func NewPool(nWorkers int) *Pool {
	if mWorkers <= 0 {
		nWorkers = 100
	}
	p := &Pool{
		q: make(chan func()error, 100),
		nWorkers: nWorkers,
		active: make(chan struct{}),
		workersDone: make(chan struct{}),
	}
	p.startPool()
	return p
}

var (
	ErrPoolClosed = errors.New("pool is shutting down")
	ErrJobPanic = errors.New("job panic was intercept")
)

func (p *Pool) startPool() {
	
	for i:=0; i < p.nWorkers; i++ {
		p.wg.Add(1)
		go func(){
			defer p.wg.Done()
			p.runWorker()
		}()
	}
	go func(){ 
		p.wg.Wait()
		close(p.workersDone) 
	}()
}

func (p *Pool) runWorker() {
	for job := range p.q {
		err := func() (e error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("worker: panic from job = %+v \n", r)
				e = ErrJobPanic
			}
		}()

		e = job()
		return e
		}()
			
		if err != nil {
			log.Printf("worker: err from job = %s \n", err.Error()) // Решил просто логировать ошибку
		}
	}
}

func (p *Pool) Submit(ctx context.Context, job func() error) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	select {
	case <-p.active:
		return ErrPoolClosed
	case <-ctx.Done():
		return ctx.Err()
	case p.q <-job:
		return nil
	}
}

func (p *Pool) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	p.once.Do(func(){
		close(p.active)
		close(p.q)
	})
	p.mu.Unlock()

	select {
	case <-p.workersDone:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}



func main(){
	pool := NewPool(100)
	for i:= 0; i < 100; i++ {
		if i%2 == 0 {
			pool.Submit(context.Background(), func() error {
				a := 0
				for j:=i; j < 200; j++{
					a = j + i + a
				}
				_ = a
				return nil
			})
		} else {
			pool.Submit(context.Background(), func() error {
				time.Sleep(time.Duration(5) * time.Second)
				return nil
			})
		}
		
	}

	pool.Shutdown(context.Background())
	ctx, cancel := 	context.WithCancel(context.Background())
	pool.Shutdown(ctx)
	cancel()
	_ = pool
}
