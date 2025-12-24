// Worker pool

//- Запускается N воркеров.
//- Задачи выполняются параллельно, но не более N одновременно.
//- Результаты выводятся в консоль по мере завершения (порядок может быть произвольным).
//- При нажатии Ctrl+C программа корректно завершает все задачи и выходит без зависаний и ошибок.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Job struct {
	Number int
}

func worker(ctx context.Context, in <-chan Job, out chan<- Job) {
	for j := range in {
		select {
		case <-ctx.Done():
			return
		default:
		}
		result := 0
		for i := 0; i < j.Number; i++ {
			result += i
		}
		select {
		case <-ctx.Done():
			return
		default:
			out <- Job{Number: result}
		}
	}
}

func WorkerPool(ctx context.Context, workers int) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	inJob := make(chan Job, workers)
	out := make(chan Job, workers)

	// some job
	go func() {
		for i := 10000; i < 1000000; i++ {
			select {
			case <-ctx.Done():
				close(inJob)
				return
			default:
			}
			inJob <- Job{Number: i}
		}
		close(inJob)
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, inJob, out)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for result := range out {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("%d \n", result.Number)
		}
	}

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	WorkerPool(ctx, 11)

}
