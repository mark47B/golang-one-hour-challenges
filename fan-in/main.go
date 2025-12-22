package main

import (
	"context"
	"sync"
)

func fanIn(ctx context.Context, chs ...<-chan string) <-chan string {
	ctxWithC, cancelF := context.WithCancel(ctx)
	out := make(chan string)
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
					out <- val
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

func generateData(id int) <-chan string {

}
