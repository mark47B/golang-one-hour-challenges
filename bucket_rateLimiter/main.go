package main

import ( 
//"runtime"
"time"
//  "fmt"
"context"
_ "log"
"errors"
// _ "net/http/pprof"
// "net/http"
	"sync"
	_ "sync/atomic"
)

type RateLimiter struct {
    rate int 
	burst int

	bucketFillSecond int
	bucket chan struct{}
	ready chan struct{}

	once sync.Once
	mu sync.RWMutex
	refillDone chan struct{}
}

var (
	ErrRateLimiterClosed = errors.New("rate limiter closed")
)

func NewRateLimiter(rate int, burst int) *RateLimiter {
	if rate <= 0 || burst <= 0 {
		return nil
	}
	rl := &RateLimiter{
		rate: rate,
		burst: burst,

		bucketFillSecond: 1,
		bucket: make(chan struct{}, burst),
		ready: make(chan struct{}),
		refillDone: make(chan struct{}),
	}

	fillInterval := time.Duration(rl.bucketFillSecond) * time.Second / time.Duration(rl.rate)
	if fillInterval == 0 {
		return nil
	}

	for i:= burst; i > 0; i--{
		rl.bucket<-struct{}{}
	}

	go func(){
		defer close(rl.refillDone)
		rl.refill() 
	}() 

	return rl
}

func (r *RateLimiter) refill() {
	ticker := time.NewTicker(time.Duration(r.bucketFillSecond) * time.Second / time.Duration(r.rate))
	defer ticker.Stop()

	for {
		select{
		case <-r.ready:
			return
		case <-ticker.C:
			select {
			case r.bucket <-struct{}{}:
			case <-r.ready:
				return
			}
		}
	}
}

func (r *RateLimiter) Allow(ctx context.Context) error {
	select {
	case <-r.bucket:
		r.mu.RLock()
		defer r.mu.RUnlock()
		select {
		case <-r.ready:
			return ErrRateLimiterClosed
		default:
			return nil
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-r.ready:
		return ErrRateLimiterClosed
	}
}

func (r *RateLimiter) Shutdown(ctx context.Context) error {
	r.mu.Lock()
	r.once.Do(func(){
		close(r.ready)
	})
	r.mu.Unlock()

	select{
	case <-r.refillDone:
		return nil
	case <- ctx.Done():
		return ctx.Err()
	}
}

func main(){
	cfgStore := NewRateLimiter(100, 20)
	_ = cfgStore
}

