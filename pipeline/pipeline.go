package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type pipelineCmd func(ctx context.Context, in <-chan int, out chan<- int)

func producer(ctx context.Context, _ <-chan int, out chan<- int) {
	for i := 0; i <= 100; i++ {
		select {
		case <-ctx.Done():
			return
		case out <- i:
		}

	}

}

func doubler(ctx context.Context, in <-chan int, out chan<- int) {

	for n := range in {
		select {
		case <-ctx.Done():
			return
		case out <- n * 2:
		}

	}
}

func printer(ctx context.Context, in <-chan int, _ chan<- int) {
	for {
		select {
		case n, ok := <-in:
			if !ok {
				return
			}
			fmt.Printf("%d-", n)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx := context.Background()
	RunPipeline(ctx, producer, doubler, printer)
}

func RunPipeline(ctx context.Context, fncs ...pipelineCmd) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(10*time.Second))
	defer cancel()

	var wg sync.WaitGroup
	out := make(chan int)
	closeCh := out

	for _, f := range fncs {
		wg.Add(1)
		in := out
		out = make(chan int)
		go func(ctx context.Context, in_ chan int, out_ chan int) {
			defer wg.Done()
			defer close(out_)
			f(ctx, in_, out_)
		}(ctx, in, out)

	}

	wg.Wait()
	close(closeCh)
}
