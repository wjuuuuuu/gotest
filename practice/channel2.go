package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// 버퍼 없는 채널 (이어 달리기)
func main() {

	track := make(chan int)
	wg.Add(1)

	go Runner(track)

	track <- 1
	wg.Wait()
}

func Runner(track chan int) {
	const maxExchanges = 4 // 최대 4회 교환
	var exchange int

	baton := <-track
	fmt.Printf("Runner %d Running with Baton\n", baton) // 현재 달리고 있는 선수

	if baton < maxExchanges {
		exchange = baton + 1
		fmt.Printf("Runner %d To The Line\n", exchange) // 새로운 선수
	}
	time.Sleep(100 * time.Millisecond)

	if baton == maxExchanges {
		fmt.Printf("Runner %d Finished, Race Over\n", baton)
		wg.Done()
		return
	}
	fmt.Printf("Runner %d Exchange With Runner %d\n", baton, exchange)
	go Runner(track)
	track <- exchange
}
