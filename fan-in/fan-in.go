package main

import (
	"context"
	"sync"
)

func fanIn(ctx context.Context, chs ...<-chan string) <-chan string {
	ctxWithC, cancelF := context.WithCancel(ctx)
	out := make(chan string, 1) // buffer
	var wg sync.WaitGroup

	wg.Add(len(chs))
	for _, ch := range chs {
		go func(ctx context.Context, in <-chan string, out chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					select { // if chan blocked
					case out <- val:
					case <-ctxWithC.Done():
						return
					}
				}
			}
		}(ctxWithC, ch, out)
	}

	go func() {
		wg.Wait()
		cancelF()
		close(out)
	}()

	return out
}
