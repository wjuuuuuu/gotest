package main

import (
	"fmt"
	"runtime"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(1)
}
func lowercase() {
	for count := 0; count < 3; count++ {
		for r := 'a'; r <= 'z'; r++ {
			fmt.Printf("%c ", r)
		}
	}
}
func uppercase() {
	for count := 0; count < 3; count++ {
		for r := 'A'; r <= 'Z'; r++ {
			fmt.Printf("%c ", r)
		}
	}
}
func main() {
	var wg sync.WaitGroup // Add : 얼마나 많은 고루틴이 있는지, Done: 일부 고루틴이 종료될 예정이므로 값을 감소, Wait: 0이 될때까지 프로그램 유지

	wg.Add(2)
	fmt.Println("Start Goroutines")

	go func() {
		lowercase()
		wg.Done()
	}()

	go func() {
		uppercase()
		wg.Done()
	}()

	fmt.Println("Waiting to finish")
	wg.Wait()

	fmt.Println("Terminating Program")
}
