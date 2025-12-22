package main

import (
	"context"
	"sync"
)

// Сделаем стратегию round-robin
func fanOut[T any](
	ctx context.Context,
	in <-chan T,
	n int,
) []<-chan T {
	ctxWithC, cancel := context.WithCancel(ctx)
	outs := make([]chan T, n)
	for i := 0; i < n; i++ {
		outs[i] = make(chan T)
	}

	outsSend := make([]chan<- T, 0, len(outs))
	outsResive := make([]<-chan T, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outsSend = append(outsSend, outs[i])
		outsResive = append(outsResive, outs[i])
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, outs []chan<- T) {
		defer wg.Done()
		n := len(outs)
		rr := 0
		for v := range in {
			select {
			case <-ctx.Done():
			case outs[rr] <- v:
			}
			rr++
			rr %= n
		}
	}(ctxWithC, outsSend)
	go func() {
		wg.Wait()
		cancel()
		for i := 0; i < len(outs); i++ {
			close(outs[i])
		}
	}()
	return outsResive
}
