package main

import "sync"

// fan in / мультиплексор каналов / мердж каналов
// берет все значения из всех входных каналов и пишет их в результирующий канал
// указатель на результирующий канал необходимо получить сразу
// после того как все входные каналы закроются - результирующий канал тоже должен быть закрыт
func merge(chans ...chan int) chan int {
	var out = make(chan int)
	var wg sync.WaitGroup
	for _, ch := range chans {
		ch := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out

}
