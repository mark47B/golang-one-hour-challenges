package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	roundTripper := &LocalFetcherRoundTripper{
		next:       http.DefaultTransport,
		logger:     os.Stdout,
		MaxRetrise: 4,
		StartDelay: time.Duration(400 * time.Millisecond),
	}
	httpClient := http.Client{
		Transport: roundTripper,
		Timeout:   3 * time.Second,
	}

	fetcher := Fetcher{
		Client: &httpClient,
	}

	crawler := Crawler{Fetcher: fetcher}

	urls := []string{"https://habr.com/ru/feed/", "https://google.com", "http://ya.ru"}

	fetcherChan := make(chan CrawlerResponse, len(urls))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := crawler.Crawl(ctx, urls, fetcherChan)
		if err != nil {
			log.Printf(err.Error())
		}
		close(fetcherChan)
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				return
			case cr, ok := <-fetcherChan:
				if !ok {
					return
				}
				fmt.Println(cr.Resp)
			}
		}
	}()

	select {
	case s := <-stop:
		log.Printf("[CATCH SIGNAL] %s\n", s.String())
		cancel()
	case <-ctx.Done():
		return
	case <-done:
	}
	wg.Wait()
}

type Crawler struct {
	Fetcher Fetcher
}

type CrawlerResponse struct {
	Resp string
	Err  error
}

func (c *Crawler) Crawl(ctx context.Context, urls []string, out chan CrawlerResponse) error {
	if ctx.Err() != nil {
		return fmt.Errorf("ctx error: %w", ctx.Err())
	}
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, url := range urls {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		go func(url string) {
			defer wg.Done()
			res, err := c.Fetcher.FetchURL(ctx, url)
			if ctx.Err() != nil {
				err = ctx.Err()
			}
			select {
			case <-ctx.Done():
			case out <- CrawlerResponse{Resp: res, Err: err}:
			}

		}(url)
	}

	wg.Wait()
	return nil
}

type LocalFetcherRoundTripper struct {
	next       http.RoundTripper
	logger     io.Writer
	MaxRetrise int
	StartDelay time.Duration
}

func (rt *LocalFetcherRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	if ctx.Err() != nil {
		return nil, fmt.Errorf("[RoundTrip] ctx error: %w", ctx.Err())
	}
	var res *http.Response
	var err error
	delay := rt.StartDelay

	for attempt := 0; attempt < rt.MaxRetrise; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		res, err = rt.next.RoundTrip(r)
		if res.StatusCode < http.StatusInternalServerError {
			break
		}
		fmt.Fprintf(rt.logger, "[FAILED REQUEST] %s", r.RequestURI)

		select {
		case <-r.Context().Done():
			return nil, r.Context().Err()
		case <-time.After(delay):
		}
		delay = delay * 2
		jitter := time.Duration(rand.Int63n(int64(delay)))
		delay = delay/2 + jitter
	}
	return res, err
}

type Fetcher struct {
	Client *http.Client
}

func (f *Fetcher) FetchURL(ctx context.Context, url string) (string, error) {
	if ctx.Err() != nil {
		return "", fmt.Errorf("ctx error: %w", ctx.Err())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error while createing request: %w", err)
	}
	r, err := f.Client.Do(req)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("request err: %s", err.Error())
		}
	}()

	if err != nil {
		return "", fmt.Errorf("get request error: %w", err)
	}

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return "", fmt.Errorf("not ok status: %s ", r.Status)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("error to read body: %w", err)
	}

	return string(b), nil
}
