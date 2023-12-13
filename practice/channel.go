package main

// 버퍼 없는 채널(공 치기)
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func player(name string, court chan int) {
	for {
		ball, ok := <-court
		if !ok {
			fmt.Printf("Player %s Won\n", name)
			return
		}
		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)
			close(court)
			return
		}
		fmt.Printf("Player %s Hit %d\n", name, ball)
		ball++
		court <- ball
	}
}
func main() {
	court := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		player("HOAN", court)
		wg.Done()
	}()
	go func() {
		player("Andrew", court)
		wg.Done()
	}()

	court <- 1
	wg.Wait()
}
