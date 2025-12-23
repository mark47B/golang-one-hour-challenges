package main

// Нужно реализовать функцию DistributedQuery, которая обращается к репликам БД и по стратегии
// first response win делает распределённый запрос. Иными словами, нужно сходить во все реплики параллельно,
// ждём пока первый запрос ответит, как первый запрос отвечает корректно, мы останавливаем все остальные запросы и отдаём ответ.

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrorNotFound = errors.New("Not Found")

type Querier interface {
	MakeQuery(ctx context.Context, host DatabaseHost) (string, error)
}

type DistributedAnswer struct {
	Resp string
	Err  error
}

// Сделаем допущение, что у нас есть экземпляр Querier
func DistributedQuery(querier Querier, replicas []DatabaseHost) (string, error) {
	if len(replicas) == 0 {
		return "", errors.New("no replicas")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Second))
	defer cancel()

	var wg sync.WaitGroup

	ch := make(chan DistributedAnswer, len(replicas))

	for i := 0; i < len(replicas); i++ {
		wg.Add(1)
		go func(ctx context.Context, host DatabaseHost) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				ch <- DistributedAnswer{Resp: "", Err: ctx.Err()}
				return
			default:
			}

			resp, err := querier.MakeQuery(ctx, host)

			select {
			case ch <- DistributedAnswer{Resp: resp, Err: err}:
			case <-ctx.Done():
				ch <- DistributedAnswer{Resp: resp, Err: ctx.Err()}
			}
		}(ctx, replicas[i])
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var lastError error
	for da := range ch {
		if da.Err == nil {
			cancel()
			wg.Wait()
			return da.Resp, nil
		}
		lastError = da.Err
	}

	return "", fmt.Errorf("all replicas failed, last error: %w", lastError)
}
