package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int // 공유 자원

func main() {
	const grs = 2

	var wg sync.WaitGroup
	wg.Add(grs)

	for i := 0; i < grs; i++ {
		go func() {
			for count := 0; count < 2; count++ {
				value := counter
				runtime.Gosched() // 다른 고루틴에게 스레드 양보하고 대기열에 들어감
				value++
				counter = value
			}
			wg.Done()
		}()

	}
	wg.Wait()
	fmt.Println("final counter", counter)

}
