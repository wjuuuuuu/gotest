package main

import (
	"fmt"
	"sync"
)

func init() {
	//runtime.GOMAXPROCS(2)
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	fmt.Println("Start Goroutines")

	go func() {
		for cnt := 0; cnt < 3; cnt++ {
			for r := 'a'; r <= 'z'; r++ {
				fmt.Printf("%c ", r)
			}
		}
		wg.Done()
	}()

	go func() {
		for cnt := 0; cnt < 3; cnt++ {
			for r := 'A'; r <= 'Z'; r++ {
				fmt.Printf("%c ", r)
			}
		}
		wg.Done()
	}()

	fmt.Println("Waiting to finish")
	wg.Wait()
	fmt.Println("Terminating Program")
}
