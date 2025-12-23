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
	if n <= 0 {
		return nil
	}
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
		rr := 0
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case outsSend[rr] <- v:
				case <-ctx.Done():
					return
				}
				rr = (rr + 1) % n
			}
		}
	}(ctx, outsSend)
	go func() {
		wg.Wait()
		for i := 0; i < len(outs); i++ {
			close(outs[i])
		}
	}()
	return outsResive
}
